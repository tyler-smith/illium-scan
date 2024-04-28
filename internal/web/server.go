package web

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/tyler-smith/iexplorer/internal/web/views"
	"github.com/tyler-smith/iexplorer/static"
)

type Server struct {
	db       *sqlx.DB
	assetDir string
	*http.ServeMux
}

func New(db *sqlx.DB, assetDir string) (*Server, error) {
	s := &Server{
		db:       db,
		assetDir: assetDir,
		ServeMux: http.NewServeMux(),
	}

	err := registerRoutes(s)
	if err != nil {
		return nil, errors.Wrap(err, "error registering routes")
	}

	return s, nil
}

func registerRoutes(s *Server) error {
	s.ServeMux.HandleFunc("/", s.Homepage)
	s.ServeMux.HandleFunc("/blocks", s.BlocksIndex)
	s.ServeMux.HandleFunc("/blocks/{id}", s.BlocksShow)
	s.ServeMux.HandleFunc("/transactions/{id}", s.TransactionsShow)

	// Load and server our static content
	staticFS, err := getEmbeddedFS()
	if err != nil {
		slog.Error("error loading static assets", "err", err)
		return errors.Wrap(err, "error loading static assets")
	}
	s.ServeMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	return nil
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

func getEmbeddedFS() (fs.FS, error) {
	return fs.Sub(static.StaticFiles, ".")
}
