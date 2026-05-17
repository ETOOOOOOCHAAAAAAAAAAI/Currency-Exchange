package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии база данных: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при открытии база данных: %w", err)
	}
	return db, nil
}
