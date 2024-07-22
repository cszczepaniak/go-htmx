package sqlite

import (
	"context"
	"testing"

	"github.com/cszczepaniak/go-htmx/internal/player/model"
	"github.com/cszczepaniak/go-htmx/internal/sql"
	"github.com/google/uuid"
	"github.com/shoenig/test"
)

func TestPlayers(t *testing.T) {
	db, err := sql.NewMemoryDB()
	test.NoError(t, err)

	ctx := context.Background()

	p := NewSQLitePlayerPersistence(db)
	test.NoError(t, p.Init(ctx))

	id1 := uuid.NewString()
	id2 := uuid.NewString()

	test.NoError(t, p.Insert(ctx, id1, "spongebob", "squarepants"))
	test.NoError(t, p.Insert(ctx, id2, "patrick", "star"))

	p1, err := p.Get(ctx, id1)
	test.NoError(t, err)
	test.Eq(
		t,
		model.Player{
			ID:        id1,
			FirstName: "spongebob",
			LastName:  "squarepants",
		},
		p1,
	)

	p2, err := p.Get(ctx, id2)
	test.NoError(t, err)
	test.Eq(
		t,
		model.Player{
			ID:        id2,
			FirstName: "patrick",
			LastName:  "star",
		},
		p2,
	)
}
