package players

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/admin/players/model"
)

type Store interface {
	Init(ctx context.Context) error

	GetPlayer(ctx context.Context, id string) (model.Player, error)
	GetPlayers(ctx context.Context, opts ...GetPlayerOpt) ([]model.Player, error)

	InsertPlayer(ctx context.Context, firstName, lastName string) (model.Player, error)
	DeletePlayer(ctx context.Context, id string) error
	AddPlayerToTeam(ctx context.Context, teamID, playerID string) error
	DeletePlayerFromTeam(ctx context.Context, teamID, playerID string) error

	GetTeam(ctx context.Context, id string) (model.Team, error)
	GetTeams(ctx context.Context, opts ...GetTeamOpt) ([]model.Team, error)

	InsertTeam(ctx context.Context) (model.Team, error)
	AddTeamToDivision(ctx context.Context, teamID, divisionID string) error
	DeleteTeamFromDivision(ctx context.Context, teamID, divisionID string) error
	DeleteTeam(ctx context.Context, id string) error
}
