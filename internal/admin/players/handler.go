package players

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/persistence"
)

func GetHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		ps, err := s.PlayerStore.GetPlayers(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, Players(ps))
	}
}

func PostHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			FirstName string `req:"form:firstName,required"`
			LastName  string `req:"form:lastName,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		_, err = s.PlayerStore.InsertPlayer(ctx, data.FirstName, data.LastName)
		if err != nil {
			return err
		}

		ps, err := s.PlayerStore.GetPlayers(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, playerList(ps))
	}
}

func DeleteHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			ID string `req:"path:id,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.PlayerStore.DeletePlayer(ctx, data.ID)
		if err != nil {
			return err
		}

		ps, err := s.PlayerStore.GetPlayers(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, playerList(ps))
	}
}
