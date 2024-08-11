package sql

import (
	"database/sql"

	"github.com/cszczepaniak/go-htmx/internal/persistence"
	divisions "github.com/cszczepaniak/go-htmx/internal/persistence/divisions/sqlite"
	players "github.com/cszczepaniak/go-htmx/internal/persistence/players/sqlite"
)

func NewPersistence(db *sql.DB) persistence.Store {
	return persistence.Store{
		PlayerStore:   players.NewSQLitePlayerPersistence(db),
		DivisionStore: divisions.NewSQLiteDivisionPersistence(db),
	}
}
