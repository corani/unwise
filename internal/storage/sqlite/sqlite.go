package sqlite

import (
	"database/sql"

	"github.com/corani/unwise/internal/config"
	_ "modernc.org/sqlite"
)

type DB struct {
	conf *config.Config
	db   *sql.DB
}

func New(conf *config.Config) (*DB, error) {
	db, err := sql.Open("sqlite", "file:quotes.db?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}

	return &DB{
		conf: conf,
		db:   db,
	}, nil
}
