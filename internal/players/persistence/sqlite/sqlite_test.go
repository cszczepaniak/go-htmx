package sqlite

import (
	"context"
	"testing"

	"github.com/cszczepaniak/go-htmx/internal/players/model"
	isql "github.com/cszczepaniak/go-htmx/internal/sql"
	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
)

func TestPlayers(t *testing.T) {
	db, err := isql.NewMemoryDB()
	test.NoError(t, err)

	ctx := context.Background()

	p := NewSQLitePlayerPersistence(db)
	test.NoError(t, p.Init(ctx))

	p1, err := p.InsertPlayer(ctx, "spongebob", "squarepants")
	test.NoError(t, err)

	p2, err := p.InsertPlayer(ctx, "patrick", "star")
	test.NoError(t, err)

	p1, err = p.GetPlayer(ctx, p1.ID)
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

	p2, err = p.GetPlayer(ctx, p2.ID)
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

	// GetPlayers should give us everything that we already validated.
	players, err := p.GetPlayers(ctx)
	must.NoError(t, err)
	test.SliceContainsAll(
		t,
		[]model.Player{p1, p2},
		players,
	)

	// DeletePlayer should work.
	must.NoError(t, p.DeletePlayer(ctx, p2.ID))

	players, err = p.GetPlayers(ctx)
	must.NoError(t, err)
	test.SliceContainsAll(
		t,
		[]model.Player{p1},
		players,
	)
}

func TestTeams(t *testing.T) {
	db, err := isql.NewMemoryDB()
	test.NoError(t, err)

	ctx := context.Background()

	p := NewSQLitePlayerPersistence(db)
	must.NoError(t, p.Init(ctx))

	p1, err := p.InsertPlayer(ctx, "spongebob", "squarepants")
	must.NoError(t, err)

	p2, err := p.InsertPlayer(ctx, "patrick", "star")
	must.NoError(t, err)

	team, err := p.InsertTeam(ctx)
	must.NoError(t, err)

	test.NotEq(t, "", team.ID)
	test.Eq(t, model.Player{}, team.Player1)
	test.Eq(t, model.Player{}, team.Player2)
	test.Eq(t, "Unnamed Team", team.Name())

	must.NoError(t, p.AddPlayerToTeam(ctx, team.ID, p1.ID))

	expP1 := p1
	expP1.TeamID = team.ID

	team, err = p.GetTeam(ctx, team.ID)
	must.NoError(t, err)
	test.NotEq(t, "", team.ID)
	test.Eq(t, expP1, team.Player1)
	test.Eq(t, model.Player{}, team.Player2)
	test.Eq(t, "squarepants", team.Name())

	test.NoError(t, p.AddPlayerToTeam(ctx, team.ID, p2.ID))

	expP2 := p2
	expP2.TeamID = team.ID

	team, err = p.GetTeam(ctx, team.ID)
	must.NoError(t, err)
	test.NotEq(t, "", team.ID)
	test.Eq(t, expP1, team.Player1)
	test.Eq(t, expP2, team.Player2)
	test.Eq(t, "squarepants/star", team.Name())

	// Adding another player at this point should error.
	p3, err := p.InsertPlayer(ctx, "anotha", "one")
	must.NoError(t, err)

	err = p.AddPlayerToTeam(ctx, team.ID, p3.ID)
	test.ErrorIs(t, err, errTeamFull)

	// Adding a player to more than one team should also error.
	t2, err := p.InsertTeam(ctx)
	must.NoError(t, err)

	err = p.AddPlayerToTeam(ctx, t2.ID, p2.ID)
	test.ErrorIs(t, err, errPlayerAlreadyOnTeam)
}
