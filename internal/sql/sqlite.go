package sql

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewMemoryDB() (*sql.DB, error) {
	return sql.Open("sqlite3", ":memory:")
}
