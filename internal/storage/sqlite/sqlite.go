package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	//this const are for print in case of possible error
	//usually operation contains func name
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	stmt, err := db.Prepare(
		`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
		`,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}
