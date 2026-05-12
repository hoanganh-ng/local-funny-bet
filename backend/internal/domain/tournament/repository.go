package tournament

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Tournament, error)
	GetByExternalCode(ctx context.Context, code string) (*Tournament, error)
	List(ctx context.Context) ([]*Tournament, error)
}
