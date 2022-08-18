package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Store struct {
	*Queries
	db *sql.DB

	config Config
}

type Config struct {
	DatabaseFilePath string
}

func NewStore(config Config) (*Store, error) {
	db, err := sql.Open("sqlite3", config.DatabaseFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "open DB")
	}
	queries := New(db)

	s := &Store{
		config:  config,
		Queries: queries,
		db:      db,
	}
	return s, nil
}
