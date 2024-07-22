package sqlite

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-htmx/internal/player/model"
)

type persistence struct {
	db *sql.DB
}

func NewSQLitePlayerPersistence(db *sql.DB) persistence {
	return persistence{
		db: db,
	}
}

func (p persistence) Init(ctx context.Context) error {
	_, err := p.db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS Players (
			ID VARCHAR(255) PRIMARY KEY,
			FirstName VARCHAR(255),
			LastName VARCHAR(255)
		)`,
	)
	return err
}

func (p persistence) Insert(ctx context.Context, id, firstName, lastName string) error {
	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO Players (ID, FirstName, LastName) VALUES (?, ?, ?)`,
		id, firstName, lastName,
	)
	return err
}

func (p persistence) Get(ctx context.Context, id string) (model.Player, error) {
	player := model.Player{
		ID: id,
	}

	err := p.db.QueryRowContext(
		ctx,
		`SELECT FirstName, LastName FROM Players WHERE ID = ?`,
		id,
	).Scan(&player.FirstName, &player.LastName)
	if err != nil {
		return model.Player{}, err
	}

	return player, nil
}
