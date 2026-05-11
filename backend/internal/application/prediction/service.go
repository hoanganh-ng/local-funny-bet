package prediction

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"wc2026/internal/domain"
	"wc2026/internal/domain/leaderboard"
	"wc2026/internal/domain/match"
	"wc2026/internal/domain/prediction"
)

type Service struct {
	predictionRepo  prediction.Repository
	matchRepo       match.Repository
	leaderboardRepo leaderboard.Repository
}

func NewService(
	predictionRepo prediction.Repository,
	matchRepo match.Repository,
	leaderboardRepo leaderboard.Repository,
) *Service {
	return &Service{
		predictionRepo:  predictionRepo,
		matchRepo:       matchRepo,
		leaderboardRepo: leaderboardRepo,
	}
}

func (s *Service) UpsertPrediction(ctx context.Context, userID, matchID, value string) error {
	m, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return fmt.Errorf("loading match: %w", err)
	}

	if m.KickoffAt.Before(time.Now()) {
		return domain.ErrLocked
	}

	var domainValue string
	switch value {
	case "home_win":
		domainValue = prediction.ValueHomeWin
	case "draw":
		domainValue = prediction.ValueDraw
	case "away_win":
		domainValue = prediction.ValueAwayWin
	default:
		return fmt.Errorf("invalid prediction value: must be home_win, draw, or away_win")
	}

	p := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userID,
		MatchID:   matchID,
		Value:     domainValue,
		UpdatedAt: time.Now(),
	}

	if err := s.predictionRepo.Upsert(ctx, p); err != nil {
		return fmt.Errorf("upserting prediction: %w", err)
	}

	return nil
}

func (s *Service) ListByMatch(ctx context.Context, leaderboardID, matchID, userID string) ([]*prediction.Prediction, error) {
	_, err := s.leaderboardRepo.GetMember(ctx, leaderboardID, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrUnauthorized
		}
		return nil, fmt.Errorf("checking membership: %w", err)
	}

	predictions, err := s.predictionRepo.ListByMatch(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("listing predictions: %w", err)
	}

	return predictions, nil
}
