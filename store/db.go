package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Store struct {
	config Config
	db     *sql.DB
}

type Config struct {
	DatabaseFilePath string
}

func NewStore(config Config) (*Store, error) {
	dsn := fmt.Sprintf("file:%s?mode=rw", config.DatabaseFilePath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "open DB")
	}
	defer db.Close()

	s := &Store{
		config: config,
		db:     db,
	}
	return s, nil
}
