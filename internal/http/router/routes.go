package router

import (
	"net/http"

	"github.com/cszczepaniak/go-htmx/internal/home"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
	"github.com/cszczepaniak/go-htmx/internal/persistence"
	"github.com/cszczepaniak/go-htmx/internal/players"
	"github.com/cszczepaniak/go-htmx/internal/web/components"
)

func Setup(p persistence.Persistence) http.Handler {
	m := http.NewServeMux()

	// Home
	httpwrap.Handle(
		m,
		"GET /",
		components.Handler(home.Home()),
	)

	// Players
	httpwrap.Handle(
		m,
		"GET /players",
		players.GetHandler(p.PlayerStore),
		components.ShellMiddleware,
	)
	httpwrap.Handle(
		m,
		"POST /players",
		players.PostHandler(p.PlayerStore),
	)
	httpwrap.Handle(
		m,
		"DELETE /players/{id}",
		players.DeleteHandler(p.PlayerStore),
	)

	return m
}
