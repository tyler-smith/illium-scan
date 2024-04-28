package web

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/tyler-smith/iexplorer/internal/db"
	"github.com/tyler-smith/iexplorer/internal/web/views"
)

func (s *Server) Homepage(w http.ResponseWriter, r *http.Request) {
	blocks, err := db.GetBlocks(s.db, 30, 0)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting blocks"))
		return
	}

	stakes, err := db.GetStakes(s.db, 100, 0)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting stakes"))
		return
	}

	proposals, err := db.GetTreasuryProposals(s.db, 100, 0)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting treasury proposals"))
		return
	}

	// Render view
	page := views.Homepage(blocks, stakes, proposals)
	render(r.Context(), w, page)
}

func (s *Server) BlocksIndex(w http.ResponseWriter, r *http.Request) {
	blocks, err := db.GetBlocks(s.db, 30, 0)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting blocks"))
		return
	}

	// Render view
	page := views.BlocksIndex(blocks)
	render(r.Context(), w, page)
}

func (s *Server) BlocksShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	block, err := db.GetBlock(s.db, id)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting block"))
		return
	}

	txs, err := db.GetTransactionsByBlockID(s.db, id)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting transactions"))
		return
	}

	// Render view
	page := views.BlocksShow(block, txs)
	render(r.Context(), w, page)
}

//func (s *Server) TransactionsIndex(w http.ResponseWriter, r *http.Request) {
//	// Render view
//	render(r.Context(), w, views.TransactionsIndex())
//}

func (s *Server) TransactionsShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tx, err := db.GetTransaction(s.db, id)
	if err != nil {
		renderError(w, errors.Wrap(err, "error getting transaction"))
		return
	}

	// Render view
	page := views.TransactionsShow(tx)
	render(r.Context(), w, page)
}
