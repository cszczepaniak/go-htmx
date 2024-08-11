package router

import (
	"io/fs"
	"net/http"

	"github.com/cszczepaniak/go-htmx/internal/admin"
	"github.com/cszczepaniak/go-htmx/internal/admin/divisions"
	"github.com/cszczepaniak/go-htmx/internal/admin/players"
	"github.com/cszczepaniak/go-htmx/internal/admin/teams"
	"github.com/cszczepaniak/go-htmx/internal/home"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/persistence"
	"github.com/cszczepaniak/go-htmx/internal/web/components"
)

func Setup(static fs.FS, s persistence.Store) http.Handler {
	m := http.NewServeMux()

	// Static files
	m.Handle(
		"GET /web/dist/",
		http.FileServerFS(static),
	)

	// Home
	httpwrap.Handle(
		m,
		"GET /",
		components.Handler(home.Home()),
	)

	httpwrap.Handle(
		m,
		"GET /admin",
		components.Handler(admin.BareAdminPage()),
	)

	httpwrap.Handle(
		m,
		"GET /admin/players",
		players.GetHandler(s),
	)
	httpwrap.Handle(
		m,
		"POST /admin/players",
		players.PostHandler(s),
	)
	httpwrap.Handle(
		m,
		"DELETE /admin/players/{id}",
		players.DeleteHandler(s),
	)

	httpwrap.Handle(
		m,
		"GET /admin/teams",
		teams.GetHandler(s),
	)
	httpwrap.Handle(
		m,
		"GET /admin/teams/list",
		teams.TeamListHandler(s),
	)
	httpwrap.Handle(
		m,
		"GET /admin/teams/{id}/edit",
		teams.EditTeamHandler(s),
	)
	httpwrap.Handle(
		m,
		"GET /admin/teams/availableplayers",
		teams.AvailablePlayersHandler(s),
	)
	httpwrap.Handle(
		m,
		"POST /admin/teams/{teamID}/player/{playerID}",
		teams.AddPlayerToTeamHandler(s),
	)
	httpwrap.Handle(
		m,
		"DELETE /admin/teams/{teamID}/player/{playerID}",
		teams.DeletePlayerFromTeamHandler(s),
	)
	httpwrap.Handle(
		m,
		"POST /admin/teams",
		teams.PostHandler(s),
	)
	httpwrap.Handle(
		m,
		"DELETE /admin/teams/{id}",
		teams.DeleteHandler(s),
	)

	httpwrap.Handle(
		m,
		"GET /admin/divisions",
		divisions.GetHandler(s),
	)
	httpwrap.Handle(
		m,
		"POST /admin/divisions",
		divisions.PostHandler(s),
	)
	httpwrap.Handle(
		m,
		"GET /admin/divisions/list",
		divisions.DivisionListHandler(s),
	)
	httpwrap.Handle(
		m,
		"GET /admin/divisions/{id}/edit",
		divisions.EditDivisionHandler(s),
	)
	httpwrap.Handle(
		m,
		"GET /admin/divisions/availableteams",
		divisions.AvailableTeamsHandler(s),
	)
	httpwrap.Handle(
		m,
		"POST /admin/divisions/{divisionID}/team/{teamID}",
		divisions.AddTeamToDivisionHandler(s),
	)
	httpwrap.Handle(
		m,
		"DELETE /admin/divisions/{divisionID}/team/{teamID}",
		divisions.DeleteTeamFromDivisionHandler(s),
	)

	return m
}
