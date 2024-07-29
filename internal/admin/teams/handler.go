package teams

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/admin/players/model"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
)

type Store interface {
	GetPlayers(ctx context.Context) ([]model.Player, error)
	InsertTeam(ctx context.Context) (model.Team, error)
	GetTeams(ctx context.Context) ([]model.Team, error)
	DeleteTeam(ctx context.Context, id string) error
}

func GetHandler(s Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return err
		}

		teams, err := s.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, Teams(ps, teams))
	}
}

func PostHandler(s Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		_, err := s.InsertTeam(ctx)
		if err != nil {
			return err
		}

		teams, err := s.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, teamList(teams))
	}
}

func DeleteHandler(s Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			ID string `req:"path:id,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.DeleteTeam(ctx, data.ID)
		if err != nil {
			return err
		}

		teams, err := s.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, teamList(teams))
	}
}
