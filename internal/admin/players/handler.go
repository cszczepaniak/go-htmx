package players

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/admin/players/model"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
)

type Store interface {
	InsertPlayer(ctx context.Context, firstName, lastName string) (model.Player, error)
	GetPlayers(ctx context.Context) ([]model.Player, error)
	DeletePlayer(ctx context.Context, id string) error
}

func GetHandler(s Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, Players(ps))
	}
}

func PostHandler(s Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			FirstName string `req:"form:firstName,required"`
			LastName  string `req:"form:lastName,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		_, err = s.InsertPlayer(ctx, data.FirstName, data.LastName)
		if err != nil {
			return err
		}

		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, playerList(ps))
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

		err = s.DeletePlayer(ctx, data.ID)
		if err != nil {
			return err
		}

		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, playerList(ps))
	}
}
