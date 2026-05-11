package match

import (
	"context"
	"encoding/json"
	"fmt"

	"wc2026/internal/domain/leaderboard"
	"wc2026/internal/domain/match"
)

type FootballClient interface {
	FetchMatches(ctx context.Context, competitionCode string) ([]*match.Match, error)
}

type Broadcaster interface {
	Broadcast(message []byte)
}

type Service struct {
	repo         match.Repository
	client       FootballClient
	leaderboards leaderboard.Repository
	broadcaster  Broadcaster
}

func NewService(repo match.Repository, client FootballClient, leaderboards leaderboard.Repository, broadcaster Broadcaster) *Service {
	return &Service{
		repo:         repo,
		client:       client,
		leaderboards: leaderboards,
		broadcaster:  broadcaster,
	}
}

func (s *Service) ListMatches(ctx context.Context, status string) ([]*match.Match, error) {
	matches, err := s.repo.List(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("listing matches: %w", err)
	}
	return matches, nil
}

func (s *Service) GetMatch(ctx context.Context, id string) (*match.Match, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting match: %w", err)
	}
	return m, nil
}

func (s *Service) FetchAndStore(ctx context.Context) error {
	matches, err := s.client.FetchMatches(ctx, "WC")
	if err != nil {
		return fmt.Errorf("fetching matches from API: %w", err)
	}

	var matchIDs []string
	for _, m := range matches {
		if err := s.repo.Upsert(ctx, m); err != nil {
			return fmt.Errorf("upserting match %s: %w", m.ID, err)
		}
		matchIDs = append(matchIDs, m.ID)
	}

	if len(matchIDs) > 0 && s.broadcaster != nil && s.leaderboards != nil {
		affectedIDs, err := s.leaderboards.GetAffectedByMatches(ctx, matchIDs)
		if err != nil {
			return fmt.Errorf("getting affected leaderboards: %w", err)
		}

		if len(affectedIDs) > 0 {
			msg := map[string]interface{}{
				"type":            "leaderboard_updated",
				"leaderboard_ids": affectedIDs,
			}
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				return fmt.Errorf("marshaling broadcast message: %w", err)
			}
			s.broadcaster.Broadcast(msgBytes)
		}
	}

	return nil
}
