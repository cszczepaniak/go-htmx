package sqlite

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-htmx/internal/admin/divisions/model"
	isql "github.com/cszczepaniak/go-htmx/internal/sql"
	"github.com/google/uuid"
)

type persistence struct {
	db *sql.DB
}

func NewSQLiteDivisionPersistence(db *sql.DB) persistence {
	return persistence{
		db: db,
	}
}

func (p persistence) Init(ctx context.Context) error {
	_, err := p.db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS Divisions (
			ID VARCHAR(255) PRIMARY KEY,
			Name VARCHAR(255)
		)`,
	)
	return err
}

func (p persistence) InsertDivision(ctx context.Context) (model.Division, error) {
	division := model.Division{
		ID: uuid.NewString(),
	}

	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO Divisions (ID) VALUES (?)`,
		division.ID,
	)
	if err != nil {
		return model.Division{}, err
	}

	return division, nil
}

func (p persistence) GetDivision(ctx context.Context, id string) (model.Division, error) {
	division := model.Division{
		ID: id,
	}

	var name sql.Null[string]
	err := p.db.QueryRowContext(
		ctx,
		`SELECT Name FROM Divisions WHERE ID = ?`,
		id,
	).Scan(&name)
	if err != nil {
		return model.Division{}, err
	}

	division.Name = name.V

	return division, nil
}

func (p persistence) GetDivisions(
	ctx context.Context,
) ([]model.Division, error) {
	q := `SELECT ID, Name FROM Divisions`

	rows, err := p.db.QueryContext(
		ctx,
		q,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var divisions []model.Division
	for rows.Next() {
		var id string
		var name sql.Null[string]
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		divisions = append(divisions, model.Division{
			ID:   id,
			Name: name.V,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return divisions, nil
}

func (p persistence) EditDivisionName(ctx context.Context, id, name string) error {
	return isql.MustExecOne(ctx, p.db, `UPDATE Divisions SET Name = ? WHERE ID = ?`, name, id)
}

func (p persistence) DeleteDivision(ctx context.Context, id string) error {
	return isql.MustExecOne(ctx, p.db, `DELETE FROM Divisions WHERE ID = ?`, id)
}
