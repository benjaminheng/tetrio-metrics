package store

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Store struct {
	config  Config
	queries *Queries
}

type Config struct {
	DatabaseFilePath string
}

func NewStore(config Config) (*Store, error) {
	db, err := sql.Open("sqlite3", config.DatabaseFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "open DB")
	}
	defer db.Close()
	queries := New(db)

	s := &Store{
		config:  config,
		queries: queries,
	}
	return s, nil
}

func (s *Store) Insert40LRecord(ctx context.Context, m *Gamemode40l) error {
	return nil
}
