package leaderboard

import "context"

type Repository interface {
	Create(ctx context.Context, lb *Leaderboard) error
	GetByID(ctx context.Context, id string) (*Leaderboard, error)
	ListByUser(ctx context.Context, userID string) ([]*Leaderboard, error)
	AddMember(ctx context.Context, m *Member) error
	GetMember(ctx context.Context, leaderboardID, userID string) (*Member, error)
	GetScores(ctx context.Context, leaderboardID string) ([]*Score, error)
	GetAffectedByMatches(ctx context.Context, matchIDs []string) ([]string, error)
}
