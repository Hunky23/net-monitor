package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("база данных: открытие соединения: %w", err)
	}
	db.SetMaxOpenConns(1)

	return db, nil
}

func Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("база данных: закрытие соединения: %w", err)
	}

	return nil
}
