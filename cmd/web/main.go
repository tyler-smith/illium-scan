package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/tyler-smith/env"

	"github.com/tyler-smith/iexplorer/internal/cmd"
	"github.com/tyler-smith/iexplorer/internal/config"
	"github.com/tyler-smith/iexplorer/internal/db"
	"github.com/tyler-smith/iexplorer/internal/web"
)

func main() {
	cmd.SetLogger()

	// Get config and connections
	conf := config.NewFromEnv()
	dbConn, err := db.NewConnection(conf.DB)
	if err != nil {
		slog.Error("error creating db connection", "err", err)
		return
	}
	defer func(dbConn *db.Connection) {
		if err := dbConn.Close(); err != nil {
			slog.Error("error closing db connection", "err", err)
			return
		}
	}(&dbConn)

	// Create server handler
	assetsDir := env.GetString("IEXP_ASSETS_DIR", "./static")
	s, err := web.New(dbConn.SQLX(), assetsDir)
	if err != nil {
		slog.Error("error creating server", "err", err)
		return
	}

	// Start server
	srv := http.Server{}
	srv.Addr = ":3000"
	srv.Handler = s
	go func() {
		slog.Info("listening on http://localhost:3000")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error listening to server", "err", err)
			return
		}
	}()

	cmd.WaitForExit()

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("error shutting down server", "err", err)
	}

	slog.Info("server shutdown")
}
