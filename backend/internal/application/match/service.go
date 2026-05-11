package match

import (
	"context"
	"fmt"

	"wc2026/internal/domain/match"
)

type FootballClient interface {
	FetchMatches(ctx context.Context, competitionCode string) ([]*match.Match, error)
}

type Service struct {
	repo   match.Repository
	client FootballClient
}

func NewService(repo match.Repository, client FootballClient) *Service {
	return &Service{
		repo:   repo,
		client: client,
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

	for _, m := range matches {
		if err := s.repo.Upsert(ctx, m); err != nil {
			return fmt.Errorf("upserting match %s: %w", m.ID, err)
		}
	}

	return nil
}
