package teams

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/persistence"
	"github.com/cszczepaniak/go-htmx/internal/persistence/players"
)

func GetHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		teams, err := s.PlayerStore.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, Teams(teams))
	}
}

func EditTeamHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			ID string `req:"path:id,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		ps, err := s.PlayerStore.GetPlayers(ctx, players.WithoutTeam())
		if err != nil {
			return err
		}

		team, err := s.PlayerStore.GetTeam(ctx, data.ID)
		if err != nil {
			return err
		}

		return req.Render(ctx, EditTeam(ps, team))
	}
}

func TeamListHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		teams, err := s.PlayerStore.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, teamList(teams))
	}
}

func AddPlayerToTeamHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			TeamID   string `req:"path:teamID,required"`
			PlayerID string `req:"path:playerID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.PlayerStore.AddPlayerToTeam(ctx, data.TeamID, data.PlayerID)
		if err != nil {
			return err
		}

		team, err := s.PlayerStore.GetTeam(ctx, data.TeamID)
		if err != nil {
			return err
		}

		req.HXTrigger("teams-damaged")
		req.HXTrigger("players-damaged")

		return req.Render(ctx, teamDetails(team))
	}
}

func DeletePlayerFromTeamHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			TeamID   string `req:"path:teamID,required"`
			PlayerID string `req:"path:playerID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		err = s.PlayerStore.DeletePlayerFromTeam(ctx, data.TeamID, data.PlayerID)
		if err != nil {
			return err
		}

		team, err := s.PlayerStore.GetTeam(ctx, data.TeamID)
		if err != nil {
			return err
		}

		req.HXTrigger("teams-damaged")
		req.HXTrigger("players-damaged")

		return req.Render(ctx, teamDetails(team))
	}
}

func AvailablePlayersHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		var data struct {
			TeamID string `req:"query:teamID,required"`
		}

		err := req.Unmarshal(&data)
		if err != nil {
			return err
		}

		ps, err := s.PlayerStore.GetPlayers(ctx, players.WithoutTeam())
		if err != nil {
			return err
		}

		return req.Render(ctx, editTeamPlayerList(data.TeamID, ps))
	}
}

func PostHandler(s persistence.Store) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		_, err := s.PlayerStore.InsertTeam(ctx)
		if err != nil {
			return err
		}

		teams, err := s.PlayerStore.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, teamList(teams))
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

		err = s.PlayerStore.DeleteTeam(ctx, data.ID)
		if err != nil {
			return err
		}

		teams, err := s.PlayerStore.GetTeams(ctx)
		if err != nil {
			return err
		}

		return req.Render(ctx, teamList(teams))
	}
}
