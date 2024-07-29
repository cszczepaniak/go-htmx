package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cszczepaniak/go-htmx/internal/admin/players/model"
	isql "github.com/cszczepaniak/go-htmx/internal/sql"
	"github.com/google/uuid"
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
			LastName VARCHAR(255),
			TeamID VARCHAR(255)
		)`,
	)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS Teams (
			ID VARCHAR(255) PRIMARY KEY
		)`,
	)
	return err
}

func (p persistence) InsertPlayer(ctx context.Context, firstName, lastName string) (model.Player, error) {
	player := model.Player{
		ID:        uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
	}

	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO Players (ID, FirstName, LastName) VALUES (?, ?, ?)`,
		player.ID, firstName, lastName,
	)
	if err != nil {
		return model.Player{}, err
	}

	return player, nil
}

func (p persistence) GetPlayer(ctx context.Context, id string) (model.Player, error) {
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

func (p persistence) GetPlayers(ctx context.Context) ([]model.Player, error) {
	rows, err := p.db.QueryContext(
		ctx,
		`SELECT ID, FirstName, LastName, TeamID FROM Players`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []model.Player
	for rows.Next() {
		var p model.Player
		var teamID sql.NullString
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &teamID)
		if err != nil {
			return nil, err
		}

		// We need to use the NullString for teamID even though we're not going to check it's valid;
		// if it's NULL we just want the empty string.
		p.TeamID = teamID.String

		players = append(players, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

func (p persistence) DeletePlayer(ctx context.Context, id string) error {
	return isql.MustExecOne(ctx, p.db, `DELETE FROM Players WHERE ID = ?`, id)
}

func (p persistence) InsertTeam(ctx context.Context) (model.Team, error) {
	id := uuid.NewString()

	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO Teams (ID) VALUES (?)`,
		id,
	)
	if err != nil {
		return model.Team{}, err
	}

	return model.Team{
		ID: id,
	}, nil
}

var (
	errTeamFull              = errors.New("team was full")
	errPlayerCouldNotBeAdded = errors.New("team not found or team was full")
	errPlayerAlreadyOnTeam   = errors.New("player was already on a team")
)

func (p persistence) AddPlayerToTeam(ctx context.Context, teamID, playerID string) (finalErr error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer isql.RollbackOnPanicOrError(&finalErr, tx)

	// Let's see if this team is already full.
	nPlayers := 0
	err = tx.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM Players WHERE TeamID = ?`,
		teamID,
	).Scan(&nPlayers)
	if err != nil {
		return err
	}
	if nPlayers >= 2 {
		return errTeamFull
	}

	err = isql.MustExecOne(
		ctx,
		tx,
		`UPDATE Players SET TeamID = ? WHERE ID = ? AND TeamID IS NULL`,
		teamID, playerID,
	)
	if err != nil {
		if errors.Is(err, isql.ErrNoRowsAffected) {
			return errPlayerAlreadyOnTeam
		}
		return err
	}

	return tx.Commit()
}

func (p persistence) GetTeam(ctx context.Context, id string) (model.Team, error) {
	t := model.Team{
		ID: id,
	}

	rows, err := p.db.QueryContext(
		ctx,
		`SELECT p.ID, p.FirstName, p.LastName
			FROM Teams t 
			LEFT JOIN Players p ON t.ID = p.TeamID
		WHERE t.ID = ?`,
		id,
	)
	if err != nil {
		return model.Team{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return model.Team{}, errors.New("team not found")
	}

	var (
		p1ID        sql.NullString
		p1FirstName sql.NullString
		p1LastName  sql.NullString
	)
	err = rows.Scan(&p1ID, &p1FirstName, &p1LastName)
	if err != nil {
		return model.Team{}, err
	}
	if !p1ID.Valid {
		// If the player fields in the first row are null, there are no players. Just return.
		return t, nil
	}

	t.Player1 = model.Player{
		ID:        p1ID.String,
		FirstName: p1FirstName.String,
		LastName:  p1LastName.String,
		TeamID:    id,
	}

	if rows.Next() {
		// There's a second player. We don't have to use nullable fields here because this can only
		// happen if the fields aren't null.
		err := rows.Scan(&t.Player2.ID, &t.Player2.FirstName, &t.Player2.LastName)
		if err != nil {
			return model.Team{}, err
		}

		t.Player2.TeamID = id
	}

	err = rows.Close()
	if err != nil {
		return model.Team{}, err
	}

	err = rows.Err()
	if err != nil {
		return model.Team{}, err
	}

	return t, nil
}
