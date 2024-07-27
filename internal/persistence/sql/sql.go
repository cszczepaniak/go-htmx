package sql

import (
	"database/sql"

	"github.com/cszczepaniak/go-htmx/internal/persistence"
	"github.com/cszczepaniak/go-htmx/internal/players/persistence/sqlite"
)

func NewPersistence(db *sql.DB) persistence.Persistence {
	return persistence.Persistence{
		PlayerStore: sqlite.NewSQLitePlayerPersistence(db),
	}
}
