package persistence

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/persistence/divisions"
	"github.com/cszczepaniak/go-htmx/internal/persistence/players"
)

type initer interface {
	Init(ctx context.Context) error
}

type Store struct {
	PlayerStore   players.Store
	DivisionStore divisions.Store
}

func (p Store) Init(ctx context.Context) error {
	for _, s := range []initer{
		p.PlayerStore,
		p.DivisionStore,
	} {
		err := s.Init(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
