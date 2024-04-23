package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/tyler-smith/iexplorer/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

func NewWriterConnection(conf config.DB) (*sql.DB, error) {
	return sql.Open(conf.Driver, conf.DSN)
}

type Connection struct {
	db   *sql.DB
	sqlx *sqlx.DB
}

func NewConnection(conf config.DB) (Connection, error) {
	db, err := sql.Open(conf.Driver, conf.DSN)
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		db:   db,
		sqlx: sqlx.NewDb(db, conf.Driver),
	}, nil
}

func (c Connection) DB() *sql.DB {
	return c.db
}

func (c Connection) SQLX() *sqlx.DB {
	return c.sqlx
}

func (c Connection) Close() error {
	err := c.db.Close()
	if err != nil {
		return err
	}

	return nil
}
