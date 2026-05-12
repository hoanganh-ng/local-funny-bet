package match

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Match, error)
	List(ctx context.Context, status string) ([]*Match, error)
	Upsert(ctx context.Context, m *Match) error
	UpsertMany(ctx context.Context, matches []*Match) error
}
