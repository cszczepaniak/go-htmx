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
		err := req.Request.ParseForm()
		if err != nil {
			return nil, err
		}

		// TODO validate the first name and the last name
		_, err = s.InsertPlayer(ctx, req.Request.FormValue("firstName"), req.Request.FormValue("lastName"))
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
		err := s.DeletePlayer(ctx, req.Request.PathValue("id"))
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
