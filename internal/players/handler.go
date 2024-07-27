package players

import (
	"context"

	"github.com/a-h/templ"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/players/model"
)

type Store interface {
	InsertPlayer(ctx context.Context, firstName, lastName string) (model.Player, error)
	GetPlayers(ctx context.Context) ([]model.Player, error)
	DeletePlayer(ctx context.Context, id string) error
}

func GetHandler(s Store) httpwrap.Handler[templ.Component] {
	return func(ctx context.Context, req httpwrap.Request) (templ.Component, error) {
		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return nil, err
		}

		return Players(ps), nil
	}
}

func PostHandler(s Store) httpwrap.Handler[templ.Component] {
	return func(ctx context.Context, req httpwrap.Request) (templ.Component, error) {
		var data struct {
			FirstName string `req:"form:firstName,required"`
			LastName  string `req:"form:lastName,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return nil, err
		}

		_, err = s.InsertPlayer(ctx, data.FirstName, data.LastName)
		if err != nil {
			return nil, err
		}

		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return nil, err
		}

		return playerList(ps), nil
	}
}

func DeleteHandler(s Store) httpwrap.Handler[templ.Component] {
	return func(ctx context.Context, req httpwrap.Request) (templ.Component, error) {
		var data struct {
			ID string `req:"path:id,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return nil, err
		}

		err = s.DeletePlayer(ctx, data.ID)
		if err != nil {
			return nil, err
		}

		ps, err := s.GetPlayers(ctx)
		if err != nil {
			return nil, err
		}

		return playerList(ps), nil
	}
}
