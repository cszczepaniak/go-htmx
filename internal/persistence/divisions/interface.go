package divisions

import (
	"context"

	"github.com/cszczepaniak/go-htmx/internal/admin/divisions/model"
)

type Store interface {
	Init(ctx context.Context) error

	GetDivision(ctx context.Context, id string) (model.Division, error)
	GetDivisions(ctx context.Context) ([]model.Division, error)

	InsertDivision(ctx context.Context) (model.Division, error)
	EditDivisionName(ctx context.Context, id, name string) error
	DeleteDivision(ctx context.Context, id string) error
}
