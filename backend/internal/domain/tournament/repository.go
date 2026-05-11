package tournament

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Tournament, error)
	List(ctx context.Context) ([]*Tournament, error)
}
