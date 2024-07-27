package components

import (
	"context"

	"github.com/a-h/templ"
	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
)

func Handler(c templ.Component) httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		return req.Render(ctx, c)
	}
}
