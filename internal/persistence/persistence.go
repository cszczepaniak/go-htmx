package persistence

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/admin/players/model"
)

type initer interface {
	Init(ctx context.Context) error
}

type PlayerStore interface {
	initer

	InsertPlayer(ctx context.Context, firstName, lastName string) (model.Player, error)
	GetPlayer(ctx context.Context, id string) (model.Player, error)
	GetPlayers(ctx context.Context) ([]model.Player, error)
	DeletePlayer(ctx context.Context, id string) error

	InsertTeam(ctx context.Context) (model.Team, error)
	AddPlayerToTeam(ctx context.Context, teamID, playerID string) error
	GetTeam(ctx context.Context, id string) (model.Team, error)
	GetTeams(ctx context.Context) ([]model.Team, error)
	DeleteTeam(ctx context.Context, id string) error
}

type Persistence struct {
	PlayerStore PlayerStore
}

func (p Persistence) Init(ctx context.Context) error {
	for _, s := range []initer{
		p.PlayerStore,
	} {
		err := s.Init(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
