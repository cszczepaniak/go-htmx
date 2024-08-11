package teams

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/persistence"
	"github.com/cszczepaniak/go-htmx/internal/persistence/players"
)

func GetHandler(s persistence.PlayerStore) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		teams, err := s.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, Teams(teams))
	}
}

func EditTeamHandler(s persistence.PlayerStore) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			ID string `req:"path:id,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		ps, err := s.GetPlayers(ctx, players.WithoutTeam())
		if err != nil {
			return err
		}

		team, err := s.GetTeam(ctx, data.ID)
		if err != nil {
			return err
		}

		return req.Render(ctx, EditTeam(ps, team))
	}
}

func TeamListHandler(s persistence.PlayerStore) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		teams, err := s.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, teamList(teams))
	}
}

func AddPlayerToTeamHandler(s persistence.PlayerStore) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			TeamID   string `req:"path:teamID,required"`
			PlayerID string `req:"path:playerID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.AddPlayerToTeam(ctx, data.TeamID, data.PlayerID)
		if err != nil {
			return err
		}

		team, err := s.GetTeam(ctx, data.TeamID)
		if err != nil {
			return err
		}

		req.HXTrigger("teams-damaged")
		req.HXTrigger("players-damaged")

		return req.Render(ctx, teamDetails(team))
	}
}

func DeletePlayerFromTeamHandler(s persistence.PlayerStore) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			TeamID   string `req:"path:teamID,required"`
			PlayerID string `req:"path:playerID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.DeletePlayerFromTeam(ctx, data.TeamID, data.PlayerID)
		if err != nil {
			return err
		}

		team, err := s.GetTeam(ctx, data.TeamID)
		if err != nil {
			return err
		}

		req.HXTrigger("teams-damaged")
		req.HXTrigger("players-damaged")

		return req.Render(ctx, teamDetails(team))
	}
}

func AvailablePlayersHandler(s persistence.PlayerStore) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			TeamID string `req:"query:teamID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		ps, err := s.GetPlayers(ctx, players.WithoutTeam())
		if err != nil {
			return err
		}

		return req.Render(ctx, editTeamPlayerList(data.TeamID, ps))
	}
}

func PostHandler(s persistence.PlayerStore) httpwrap.Handler {
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

func DeleteHandler(s persistence.PlayerStore) httpwrap.Handler {
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
