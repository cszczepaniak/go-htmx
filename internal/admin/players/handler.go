package players

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/persistence"
)

func GetHandler(s persistence.PlayerStore) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, Players(ps))
	}
}

func PostHandler(s persistence.PlayerStore) httpwrap.Handler {
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

func DeleteHandler(s persistence.PlayerStore) httpwrap.Handler {
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
