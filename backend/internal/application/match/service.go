package match

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"wc2026/internal/domain/leaderboard"
	"wc2026/internal/domain/match"
	"wc2026/internal/domain/tournament"
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
	tournaments  tournament.Repository
	leaderboards leaderboard.Repository
	broadcaster  Broadcaster
}

func NewService(repo match.Repository, client FootballClient, tournaments tournament.Repository, leaderboards leaderboard.Repository, broadcaster Broadcaster) *Service {
	return &Service{
		repo:         repo,
		client:       client,
		tournaments:  tournaments,
		leaderboards: leaderboards,
		broadcaster:  broadcaster,
	}
}

func (s *Service) ListMatches(ctx context.Context, status string, limit int, after time.Time) ([]*match.Match, error) {
	matches, err := s.repo.List(ctx, status, limit, after)
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
	active, err := s.tournaments.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("getting active tournament: %w", err)
	}

	if active.ExternalID == nil {
		return fmt.Errorf("active tournament has no external_id")
	}

	matches, err := s.client.FetchMatches(ctx, *active.ExternalID)
	if err != nil {
		return fmt.Errorf("fetching matches from API: %w", err)
	}

	for _, m := range matches {
		m.TournamentID = active.ID
	}

	if err := s.repo.UpsertMany(ctx, matches); err != nil {
		return fmt.Errorf("upserting matches: %w", err)
	}

	if len(matches) > 0 && s.broadcaster != nil && s.leaderboards != nil {
		matchIDs := make([]string, len(matches))
		for i, m := range matches {
			matchIDs[i] = m.ID
		}

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
