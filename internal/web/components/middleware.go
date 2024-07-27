package components

import (
	"context"

	"github.com/a-h/templ"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
)

func Handler(c templ.Component) httpwrap.Handler[templ.Component] {
	return func(context.Context, httpwrap.Request) (templ.Component, error) {
		return c, nil
	}
}

func ShellMiddleware(next httpwrap.Handler[templ.Component]) httpwrap.Handler[templ.Component] {
	return func(ctx context.Context, req httpwrap.Request) (templ.Component, error) {
		c, err := next(ctx, req)
		if err != nil {
			return nil, err
		}

		return Head(c), nil
	}
}
