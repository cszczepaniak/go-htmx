package sqlite

import (
	"context"
	"testing"

	"github.com/cszczepaniak/go-htmx/internal/admin/players/model"
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

	team1, err := p.InsertTeam(ctx)
	must.NoError(t, err)

	test.NotEq(t, "", team1.ID)
	test.Eq(t, model.Player{}, team1.Player1)
	test.Eq(t, model.Player{}, team1.Player2)
	test.Eq(t, "Unnamed Team", team1.Name())

	must.NoError(t, p.AddPlayerToTeam(ctx, team1.ID, p1.ID))

	expP1 := p1
	expP1.TeamID = team1.ID

	team1, err = p.GetTeam(ctx, team1.ID)
	must.NoError(t, err)
	test.NotEq(t, "", team1.ID)
	test.Eq(t, expP1, team1.Player1)
	test.Eq(t, model.Player{}, team1.Player2)
	test.Eq(t, "squarepants", team1.Name())

	test.NoError(t, p.AddPlayerToTeam(ctx, team1.ID, p2.ID))

	expP2 := p2
	expP2.TeamID = team1.ID

	team1, err = p.GetTeam(ctx, team1.ID)
	must.NoError(t, err)
	test.NotEq(t, "", team1.ID)
	test.Eq(t, expP1, team1.Player1)
	test.Eq(t, expP2, team1.Player2)
	test.Eq(t, "squarepants/star", team1.Name())

	// Adding another player at this point should error.
	p3, err := p.InsertPlayer(ctx, "anotha", "one")
	must.NoError(t, err)

	err = p.AddPlayerToTeam(ctx, team1.ID, p3.ID)
	test.ErrorIs(t, err, errTeamFull)

	// Adding a player to more than one team should also error.
	team2, err := p.InsertTeam(ctx)
	must.NoError(t, err)

	err = p.AddPlayerToTeam(ctx, team2.ID, p2.ID)
	test.ErrorIs(t, err, errPlayerAlreadyOnTeam)

	// Set up two more teams so we can test GetTeams.
	must.NoError(t, p.AddPlayerToTeam(ctx, team2.ID, p3.ID))
	team2, err = p.GetTeam(ctx, team2.ID)
	must.NoError(t, err)

	team3, err := p.InsertTeam(ctx)
	must.NoError(t, err)

	teams, err := p.GetTeams(ctx)
	must.NoError(t, err)

	teamsByID := make(map[string]model.Team, len(teams))
	for _, team := range teams {
		teamsByID[team.ID] = team
	}

	must.MapContainsKey(t, teamsByID, team1.ID)
	test.Eq(t, team1, teamsByID[team1.ID])

	must.MapContainsKey(t, teamsByID, team2.ID)
	test.Eq(t, team2, teamsByID[team2.ID])

	must.MapContainsKey(t, teamsByID, team3.ID)
	test.Eq(t, team3, teamsByID[team3.ID])
}
