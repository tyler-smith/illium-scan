package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/project-illium/ilxd/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/tyler-smith/iexplorer/internal/cmd"
	"github.com/tyler-smith/iexplorer/internal/config"
	"github.com/tyler-smith/iexplorer/internal/db"
)

type doneCh chan struct{}
type quitCh chan doneCh

func main() {
	cmd.SetLogger()

	// Get config and connections
	conf := config.NewFromEnv()
	dbConn, err := db.NewWriterConnection(conf.DB)
	if err != nil {
		slog.Error("error connecting to database", "err", err)
		return
	}
	defer func(dbConn *sql.DB) {
		if err := dbConn.Close(); err != nil {
			slog.Error("error closing database connection", "err", err)
			return
		}
	}(dbConn)

	// Prepare statements.
	stmts, err := db.NewCachedWriterStmts(dbConn)
	if err != nil {
		slog.Error("error preparing statements", "err", err)
		return
	}

	// Run the indexer until we're asked to quit
	quit := start(conf.Indexer, stmts, dbConn)
	cmd.WaitForExit()

	// Stop the indexer
	done := make(doneCh)
	quit <- done
	<-done
}

func start(conf config.Indexer, stmts db.CachedWriterStmts, sqlDB *sql.DB) quitCh {
	quit := make(quitCh)
	errorCount := 0
	go func() {
		for {
			// Stop if we've been asked to quit
			select {
			case done := <-quit:
				close(done)
				return
			default:
			}

			// Connect to the server and wait for notifications
			err := connectAndWatch(conf, sqlDB, stmts)

			// If we have an error log in, then wait before reconnecting
			if err != nil {
				slog.Error("error watching for changes", "err", err)

				errorCount += 1
				select {
				case <-time.After(backoffSeconds(errorCount)):
					continue
				case done := <-quit:
					close(done)
					return
				}
			}

			errorCount = 0
		}
	}()
	return quit
}

func connectAndWatch(conf config.Indexer, sqlDB *sql.DB, stmts db.CachedWriterStmts) error {
	// Connect to the gRPC server
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
		//grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(conf.GRPCServerAddr, grpcOpts...)
	if err != nil {
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error("error closing grpc connection", "err", err)
			return
		}
	}(conn)
	client := pb.NewBlockchainServiceClient(conn)

	// Create new block subscription.
	blockSub, err := client.SubscribeBlocks(context.Background(), &pb.SubscribeBlocksRequest{
		FullBlock:        true,
		FullTransactions: true,
	})
	if err != nil {
		return err
	}

	// Handle block notifications until stopped.
	for {
		slog.Info("Waiting for block")
		notification, err := blockSub.Recv()
		if err != nil {
			return err
		}
		err = handleBlockNotification(sqlDB, stmts, notification)
		if err != nil {
			return err
		}
	}
}

func handleBlockNotification(sqlDB *sql.DB, stmts db.CachedWriterStmts, notification *pb.BlockNotification) error {
	header := notification.GetBlockInfo()
	slog.Info("Received block", "height", header.GetHeight())

	// Create DB transaction
	ctx, cancelFn := newDBCtx()
	defer cancelFn()
	dbTx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		err := dbTx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("error rolling back transaction", "err", err)
			return
		}
	}()

	// Prepare statements for the transaction.
	stmts = stmts.ForTx(dbTx)

	// Insert each transaction.
	for _, tx := range notification.GetTransactions() {
		if err := db.InsertTransaction(stmts, header.GetBlock_ID(), tx.GetTransaction()); err != nil {
			// Seeing the same transaction is not an error.
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
				slog.Debug("transaction already exists", "id", formatID(tx.GetTransaction().ID().Bytes()))
				continue
			}

			return err
		}
	}

	// Insert the block.
	if err := db.InsertBlock(stmts, header); err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			slog.Debug("block already exists", "id", hex.EncodeToString(header.GetBlock_ID()))
			return nil
		}
		return err
	}

	// Commit the DB transaction.
	if err := dbTx.Commit(); err != nil {
		return err
	}

	slog.Info("Inserted block", "height", header.GetHeight())
	return nil
}

func newDBCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func backoffSeconds(errors int) time.Duration {
	backoff := errors * 2
	if backoff > 60 {
		backoff = 60
	}
	return time.Duration(backoff) * time.Second
}

func formatID(id []byte) string {
	return hex.EncodeToString(id)
}
