package leaderboard

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"wc2026/internal/domain"
	"wc2026/internal/domain/auth"
	"wc2026/internal/domain/leaderboard"
)

type Service struct {
	leaderboardRepo leaderboard.Repository
	inviteSecret    string
}

func NewService(leaderboardRepo leaderboard.Repository, inviteSecret string) *Service {
	return &Service{
		leaderboardRepo: leaderboardRepo,
		inviteSecret:    inviteSecret,
	}
}

func (s *Service) CreateLeaderboard(ctx context.Context, userID, name string) (*leaderboard.Leaderboard, error) {
	lb := &leaderboard.Leaderboard{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	if err := s.leaderboardRepo.Create(ctx, lb); err != nil {
		return nil, fmt.Errorf("creating leaderboard: %w", err)
	}

	m := &leaderboard.Member{
		ID:            uuid.New().String(),
		LeaderboardID: lb.ID,
		UserID:        userID,
		Role:          leaderboard.RoleOwner,
		JoinedAt:      time.Now(),
	}

	if err := s.leaderboardRepo.AddMember(ctx, m); err != nil {
		return nil, fmt.Errorf("adding owner as member: %w", err)
	}

	return lb, nil
}

func (s *Service) GetLeaderboard(ctx context.Context, id string) (*leaderboard.Leaderboard, error) {
	lb, err := s.leaderboardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting leaderboard: %w", err)
	}
	return lb, nil
}

func (s *Service) ListMyLeaderboards(ctx context.Context, userID string) ([]*leaderboard.Leaderboard, error) {
	leaderboards, err := s.leaderboardRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing leaderboards: %w", err)
	}
	return leaderboards, nil
}

func (s *Service) GenerateInvite(ctx context.Context, leaderboardID, userID string) (string, time.Time, error) {
	member, err := s.leaderboardRepo.GetMember(ctx, leaderboardID, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			return "", time.Time{}, domain.ErrUnauthorized
		}
		return "", time.Time{}, fmt.Errorf("getting member: %w", err)
	}

	if member.Role != leaderboard.RoleOwner {
		return "", time.Time{}, domain.ErrUnauthorized
	}

	token := auth.SignInviteToken(leaderboardID, s.inviteSecret)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	return token, expiresAt, nil
}

func (s *Service) JoinLeaderboard(ctx context.Context, userID, inviteToken string) error {
	leaderboardID, err := auth.VerifyInviteToken(inviteToken, s.inviteSecret)
	if err != nil {
		return err
	}

	m := &leaderboard.Member{
		ID:            uuid.New().String(),
		LeaderboardID: leaderboardID,
		UserID:        userID,
		Role:          leaderboard.RoleMember,
		JoinedAt:      time.Now(),
	}

	if err := s.leaderboardRepo.AddMember(ctx, m); err != nil {
		return fmt.Errorf("adding member: %w", err)
	}

	return nil
}
