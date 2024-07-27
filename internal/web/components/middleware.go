package components

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
)

func Handler(c templ.Component) httpwrap.Handler[templ.Component] {
	return func(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
		return c, nil
	}
}

func ShellMiddleware(next httpwrap.Handler[templ.Component]) httpwrap.Handler[templ.Component] {
	return func(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
		c, err := next(w, r)
		if err != nil {
			return nil, err
		}

		return Head(c), nil
	}
}
