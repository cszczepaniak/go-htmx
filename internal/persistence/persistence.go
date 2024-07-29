package persistence

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/admin/players"
)

type initer interface {
	Init(ctx context.Context) error
}

type Persistence struct {
	PlayerStore players.Store
}

func (p Persistence) Init(ctx context.Context) error {
	for _, s := range []any{
		p.PlayerStore,
	} {
		i, ok := s.(initer)
		if !ok {
			continue
		}

		err := i.Init(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
