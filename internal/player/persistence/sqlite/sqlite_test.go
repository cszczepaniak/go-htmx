package sqlite

import (
	"context"
	"testing"

	"github.com/cszczepaniak/go-htmx/internal/player/model"
	"github.com/cszczepaniak/go-htmx/internal/sql"
	"github.com/shoenig/test"
)

func TestPlayers(t *testing.T) {
	db, err := sql.NewMemoryDB()
	test.NoError(t, err)

	ctx := context.Background()

	p := NewSQLitePlayerPersistence(db)
	test.NoError(t, p.Init(ctx))

	p1, err := p.Insert(ctx, "spongebob", "squarepants")
	test.NoError(t, err)

	p2, err := p.Insert(ctx, "patrick", "star")
	test.NoError(t, err)

	p1, err = p.Get(ctx, p1.ID)
	test.NoError(t, err)
	test.Eq(
		t,
		model.Player{
			ID:        p1.ID,
			FirstName: "spongebob",
			LastName:  "squarepants",
		},
		p1,
	)

	p2, err = p.Get(ctx, p2.ID)
	test.NoError(t, err)
	test.Eq(
		t,
		model.Player{
			ID:        p2.ID,
			FirstName: "patrick",
			LastName:  "star",
		},
		p2,
	)
}
