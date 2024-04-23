package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/project-illium/ilxd/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/tyler-smith/iexplorer/internal/cmd"
	"github.com/tyler-smith/iexplorer/internal/config"
	"github.com/tyler-smith/iexplorer/internal/db"
)

type signalCh chan struct{}

func main() {
	cmd.SetLogger()

	// Get config and connections
	conf := config.NewFromEnv()
	sqlDB, err := db.NewWriterConnection(conf.DB)
	if err != nil {
		slog.Error("error connecting to database", "err", err)
		return
	}
	defer func(dbConn *sql.DB) {
		if err := dbConn.Close(); err != nil {
			slog.Error("error closing database connection", "err", err)
			return
		}
	}(sqlDB)

	// Prepare statements.
	stmts, err := db.NewCachedWriterStmts(sqlDB)
	if err != nil {
		slog.Error("error preparing statements", "err", err)
		return
	}

	start(conf.Indexer, stmts, sqlDB)

	cmd.WaitForExit()
}

func start(conf config.Indexer, stmts db.CachedWriterStmts, sqlDB *sql.DB) signalCh {
	quitCh := make(signalCh)
	go func() {
		for {
			err := connectAndWatch(conf, sqlDB, stmts)
			if err != nil {
				slog.Error("error watching for changes", "err", err)
			}
		}
	}()
	return quitCh
}

func connectAndWatch(conf config.Indexer, sqlDB *sql.DB, stmts db.CachedWriterStmts) error {
	// Connect to the gRPC server
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
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
		if !errors.Is(err, sql.ErrTxDone) {
			slog.Error("error rolling back transaction", "err", err)
			return
		}
	}()

	// Prepare statements for the transaction.
	stmts = stmts.ForTx(dbTx)

	// Insert each transaction.
	for _, tx := range notification.GetTransactions() {
		if err := db.InsertTransaction(stmts, header.GetBlock_ID(), tx.GetTransaction()); err != nil {
			return err
		}
	}

	// Insert the block.
	if err := db.InsertBlock(stmts, header); err != nil {
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
