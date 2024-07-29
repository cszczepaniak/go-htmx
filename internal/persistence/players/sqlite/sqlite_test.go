package sqlite

import (
	"context"
	"fmt"
	"testing"

	"github.com/cszczepaniak/go-htmx/internal/admin/players/model"
	isql "github.com/cszczepaniak/go-htmx/internal/sql"
	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
)

type playerServiceTester struct {
	persistence
}

func newPlayerServiceTester(t testing.TB) playerServiceTester {
	db, err := isql.NewMemoryDB()
	test.NoError(t, err)

	ctx := context.Background()

	p := NewSQLitePlayerPersistence(db)
	test.NoError(t, p.Init(ctx))

	return playerServiceTester{
		persistence: p,
	}
}

func (p playerServiceTester) seedPlayers(
	t testing.TB,
	n int,
) []model.Player {
	t.Helper()

	res := make([]model.Player, 0, n)
	for i := range n {
		p, err := p.InsertPlayer(
			context.Background(),
			fmt.Sprintf("first%d", i),
			fmt.Sprintf("last%d", i),
		)
		must.NoError(t, err)

		res = append(res, p)
	}

	return res
}

func (p playerServiceTester) seedPlayer(
	t testing.TB,
) model.Player {
	t.Helper()

	return p.seedPlayers(t, 1)[0]
}

func (p playerServiceTester) getPlayer(t testing.TB, id string) model.Player {
	t.Helper()

	player, err := p.GetPlayer(context.Background(), id)
	must.NoError(t, err)

	return player
}

func (p playerServiceTester) getPlayers(t testing.TB) map[string]model.Player {
	t.Helper()

	players, err := p.GetPlayers(context.Background())
	must.NoError(t, err)

	byID := make(map[string]model.Player)
	for _, player := range players {
		byID[player.ID] = player
	}

	return byID
}

func (p playerServiceTester) seedTeams(t testing.TB, n int) []model.Team {
	t.Helper()

	teams := make([]model.Team, 0, n)
	for range n {
		team, err := p.InsertTeam(context.Background())
		must.NoError(t, err)

		teams = append(teams, team)
	}

	return teams
}

func (p playerServiceTester) seedTeam(t testing.TB) model.Team {
	t.Helper()
	return p.seedTeams(t, 1)[0]
}

func (p playerServiceTester) addPlayerToTeam(
	t testing.TB,
	player model.Player,
	team model.Team,
) {
	t.Helper()

	must.NoError(t, p.AddPlayerToTeam(context.Background(), team.ID, player.ID))
}

func TestInsertPlayers(t *testing.T) {
	p := newPlayerServiceTester(t)
	ctx := context.Background()

	p1, err := p.InsertPlayer(ctx, "spongebob", "squarepants")
	test.NoError(t, err)

	p2, err := p.InsertPlayer(ctx, "patrick", "star")
	test.NoError(t, err)

	players := p.getPlayers(t)

	must.MapContainsKey(t, players, p1.ID)
	test.Eq(t, players[p1.ID], p1)

	must.MapContainsKey(t, players, p2.ID)
	test.Eq(t, players[p2.ID], p2)
}

func TestGetPlayers(t *testing.T) {
	p := newPlayerServiceTester(t)
	ctx := context.Background()

	players := p.seedPlayers(t, 3)

	got, err := p.GetPlayers(ctx)
	must.NoError(t, err)
	test.SliceContainsAll(
		t,
		got,
		players,
	)

	// Add one of the players to a team.
	team := p.seedTeam(t)
	p.addPlayerToTeam(t, players[1], team)

	// The player with a team should not be returned.
	got, err = p.GetPlayers(ctx, WithoutTeam())
	must.NoError(t, err)
	test.SliceContainsAll(
		t,
		got,
		[]model.Player{players[0], players[2]},
	)
}

func TestDeletePlayer(t *testing.T) {
	p := newPlayerServiceTester(t)
	ctx := context.Background()

	players := p.seedPlayers(t, 3)

	must.NoError(t, p.DeletePlayer(ctx, players[1].ID))

	got, err := p.GetPlayers(ctx)
	must.NoError(t, err)
	test.SliceContainsAll(
		t,
		got,
		[]model.Player{players[0], players[2]},
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
