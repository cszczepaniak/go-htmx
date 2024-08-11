package divisions

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/persistence"
	"github.com/cszczepaniak/go-htmx/internal/persistence/players"
)

func GetHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		ds, err := s.DivisionStore.GetDivisions(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, divisions(ds))
	}
}

func EditDivisionHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			ID string `req:"path:id,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		d, err := s.DivisionStore.GetDivision(ctx, data.ID)
		if err != nil {
			return err
		}

		teamsWithoutDivision, err := s.PlayerStore.GetTeams(ctx, players.WithoutDivision())
		if err != nil {
			return err
		}

		teamsOnDivision, err := s.PlayerStore.GetTeams(ctx, players.InDivision(data.ID))
		if err != nil {
			return err
		}

		return req.Render(ctx, EditDivision(teamsWithoutDivision, d, teamsOnDivision))
	}
}

func DivisionListHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		teams, err := s.DivisionStore.GetDivisions(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, divisionList(teams))
	}
}

func AddTeamToDivisionHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			DivisionID string `req:"path:divisionID,required"`
			TeamID     string `req:"path:teamID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.PlayerStore.AddTeamToDivision(ctx, data.TeamID, data.DivisionID)
		if err != nil {
			return err
		}

		division, err := s.DivisionStore.GetDivision(ctx, data.DivisionID)
		if err != nil {
			return err
		}

		teamsOnDivision, err := s.PlayerStore.GetTeams(ctx, players.InDivision(division.ID))
		if err != nil {
			return err
		}

		req.HXTrigger("divisions-damaged")
		req.HXTrigger("teams-damaged")

		return req.Render(ctx, divisionDetails(division, teamsOnDivision))
	}
}

func DeleteTeamFromDivisionHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			DivisionID string `req:"path:divisionID,required"`
			TeamID     string `req:"path:teamID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.PlayerStore.DeleteTeamFromDivision(ctx, data.TeamID, data.DivisionID)
		if err != nil {
			return err
		}

		division, err := s.DivisionStore.GetDivision(ctx, data.DivisionID)
		if err != nil {
			return err
		}

		teams, err := s.PlayerStore.GetTeams(ctx, players.InDivision(division.ID))
		if err != nil {
			return err
		}

		req.HXTrigger("teams-damaged")
		req.HXTrigger("divisions-damaged")

		return req.Render(ctx, divisionDetails(division, teams))
	}
}

func AvailableTeamsHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			DivisionID string `req:"query:divisionID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		ts, err := s.PlayerStore.GetTeams(ctx, players.WithoutDivision())
		if err != nil {
			return err
		}

		return req.Render(ctx, editDivisionTeamList(data.DivisionID, ts))
	}
}

func PostHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		_, err := s.DivisionStore.InsertDivision(ctx)
		if err != nil {
			return err
		}

		ds, err := s.DivisionStore.GetDivisions(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, divisionList(ds))
	}
}
