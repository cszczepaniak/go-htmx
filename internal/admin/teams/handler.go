package teams

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/http/httpwrap"
)

func GetHandler() httpwrap.Handler {
	return func(ctx context.Context, req httpwrap.Request) error {
		return req.Render(ctx, Teams())
	}
}
