package store

import (
	"database/sql"
	"log"
	"os"

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
	db, err := sql.Open("sqlite3", config.DatabaseFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "open DB")
	}
	defer db.Close()

	s := &Store{
		config: config,
		db:     db,
	}
	if _, err := os.Stat(config.DatabaseFilePath); errors.Is(err, os.ErrNotExist) {
		log.Printf("initializing database %s\n", config.DatabaseFilePath)
		err := s.initTables()
		if err != nil {
			return nil, errors.Wrap(err, "init database tables")
		}

	}
	return s, nil
}

func (s *Store) initTables() error {
	stmt := `
CREATE TABLE IF NOT EXISTS gamemode_40l (
  id integer not null primary key,
  played_at timestamp,
  time double,
  finesse_percent double,
  finesse_faults integer,
  total_pieces integer,
  rawData text
);`
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
