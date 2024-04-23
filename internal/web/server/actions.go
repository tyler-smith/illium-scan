package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"

	"github.com/tyler-smith/iexplorer/internal/db"
	"github.com/tyler-smith/iexplorer/internal/web/viewmodels"
	"github.com/tyler-smith/iexplorer/internal/web/views"
)

func (s *Server) BlocksIndex(w http.ResponseWriter, r *http.Request) {
	blocks, err := db.GetBlocks(s.dbConn.SQLX(), 30, 0)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting blocks"))
	}

	// Render view
	viewModel := viewmodels.NewBlocksIndex(blocks)
	page := views.BlocksIndex(viewModel)
	render(r.Context(), w, page)
}

func (s *Server) BlocksShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	block, err := db.GetBlock(s.dbConn.SQLX(), id)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting block"))
	}

	// Render view
	viewModel := viewmodels.NewBlocksShow(block)
	page := views.BlocksShow(viewModel)
	render(r.Context(), w, page)
}

//func (s *Server) TransactionsIndex(w http.ResponseWriter, r *http.Request) {
//	// Render view
//	render(r.Context(), w, views.TransactionsIndex())
//}

func (s *Server) TransactionsShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	tx, err := db.GetTransaction(s.dbConn.SQLX(), id)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting transaction"))
	}

	// debug
	txJSON, err := json.Marshal(tx)
	if err != nil {
		renderError(w, errors.Wrap(err, "error marshalling transaction"))
	}
	slog.Debug("transaction", "tx", string(txJSON))

	// Render view
	viewModel := viewmodels.NewTransactionsShow(tx)
	page := views.TransactionsShow(viewModel)
	render(r.Context(), w, page)
}
