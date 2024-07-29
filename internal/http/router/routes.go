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

func Setup(static fs.FS, p persistence.Persistence) http.Handler {
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
		players.GetHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"POST /admin/players",
		players.PostHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"DELETE /admin/players/{id}",
		players.DeleteHandler(p.PlayerStore),
	)

	httpwrap.Handle(
		m,
		"GET /admin/teams",
		teams.GetHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"GET /admin/teams/list",
		teams.TeamListHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"GET /admin/teams/{id}/edit",
		teams.EditTeamHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"POST /admin/teams/{teamID}/player/{playerID}",
		teams.AddPlayerToTeamHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"POST /admin/teams",
		teams.PostHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"DELETE /admin/teams/{id}",
		teams.DeleteHandler(p.PlayerStore),
	)

	httpwrap.Handle(
		m,
		"GET /admin/divisions",
		divisions.GetHandler(),
	)

	// Players

	return m
}
