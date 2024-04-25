package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/pkg/errors"

	"github.com/tyler-smith/iexplorer/internal/db"
	"github.com/tyler-smith/iexplorer/internal/web/views"
)

type Server struct {
	dbConn   db.Connection
	assetDir string
	*http.ServeMux
}

func New(dbConn db.Connection, assetDir string) *Server {
	s := &Server{
		dbConn:   dbConn,
		assetDir: assetDir,
		ServeMux: http.NewServeMux(),
	}

	registerRoutes(s)

	return s
}

func registerRoutes(s *Server) {
	s.ServeMux.HandleFunc("/", s.Homepage)
	s.ServeMux.HandleFunc("/blocks", s.BlocksIndex)
	s.ServeMux.HandleFunc("/blocks/{id}", s.BlocksShow)
	s.ServeMux.HandleFunc("/transactions/{id}", s.TransactionsShow)
	s.ServeMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.assetDir))))
	//http.Handle("/", renderView(templates.Homepage))
	//http.Handle("/blocks/{id}", renderView(templates.BlockShow))
	//http.Handle("/transactions/{id}", renderView(templates.TxShow))
	//http.Handle("/validators", renderView(templates.ValidatorIndex))
	//http.Handle("/validators/{id}", renderView(templates.ValidatorShow))
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))
}

func renderError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	slog.Error("server error", "err", err)
}

func render(ctx context.Context, w http.ResponseWriter, page templ.Component) {
	w.Header().Set("Content-Type", "text/html")
	layout := views.Layout(page)
	if err := layout.Render(ctx, w); err != nil {
		renderError(w, errors.Wrap(err, "error rendering page"))
	}
}
