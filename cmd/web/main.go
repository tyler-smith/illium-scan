package main

import (
	"log/slog"
	"net/http"

	"github.com/tyler-smith/env"

	"github.com/tyler-smith/iexplorer/internal/cmd"
	"github.com/tyler-smith/iexplorer/internal/config"
	"github.com/tyler-smith/iexplorer/internal/db"
	"github.com/tyler-smith/iexplorer/internal/web/server"
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

	// Create server
	assetsDir := env.GetString("IEXP_ASSETS_DIR", "./static")
	s := server.New(dbConn, assetsDir)
	go func() {
		slog.Info("listening on http://localhost:3000")
		err := http.ListenAndServe(":3000", s)
		if err != nil {
			slog.Error("error listening to server", "err", err)
			return
		}
	}()

	cmd.WaitForExit()
}
