package prediction

import "context"

type Repository interface {
	Upsert(ctx context.Context, p *Prediction) error
	GetByUserAndMatch(ctx context.Context, userID, matchID string) (*Prediction, error)
	ListByMatch(ctx context.Context, matchID string) ([]*Prediction, error)
	ListByUser(ctx context.Context, userID string) ([]*PredictionWithMatch, error)
}
