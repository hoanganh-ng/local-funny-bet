package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"wc2026/internal/domain/auth"
	"wc2026/internal/domain/user"
)

type Service struct {
	userRepo         user.Repository
	refreshTokenRepo auth.RefreshTokenRepository
	provider         auth.Provider
	tokenSigner      auth.TokenSigner
	refreshSecret    string
}

func NewService(userRepo user.Repository, refreshTokenRepo auth.RefreshTokenRepository, provider auth.Provider, tokenSigner auth.TokenSigner, refreshSecret string) *Service {
	return &Service{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		provider:         provider,
		tokenSigner:      tokenSigner,
		refreshSecret:    refreshSecret,
	}
}

func (s *Service) GetAuthURL(state string) string {
	return s.provider.GetAuthURL(state)
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (id, email, name string, avatarURL *string, err error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", "", nil, err
	}
	return u.ID, u.Email, u.Name, u.AvatarURL, nil
}

func (s *Service) HandleGoogleCallback(ctx context.Context, code string) (accessToken, refreshTokenStr string, err error) {
	userInfo, err := s.provider.VerifyToken(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("verifying token: %w", err)
	}

	now := time.Now()
	u := &user.User{
		ID:        uuid.New().String(),
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: &userInfo.AvatarURL,
		CreatedAt: now,
	}

	if err := s.userRepo.Upsert(ctx, u); err != nil {
		return "", "", fmt.Errorf("upserting user: %w", err)
	}

	accessToken, err = s.tokenSigner.Sign(u.ID, u.Email, u.Name, s.provider.Name(), 15*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("signing access token: %w", err)
	}

	refreshTokenStr, err = s.createRefreshToken(u.ID)
	if err != nil {
		return "", "", fmt.Errorf("creating refresh token: %w", err)
	}

	return accessToken, refreshTokenStr, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshTokenStr string) (string, error) {
	rt, err := s.refreshTokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = s.refreshTokenRepo.Delete(ctx, refreshTokenStr)
		return "", fmt.Errorf("refresh token expired: %w", fmt.Errorf("unauthorized"))
	}

	u, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return "", fmt.Errorf("getting user: %w", err)
	}

	accessToken, err := s.tokenSigner.Sign(u.ID, u.Email, u.Name, "google", 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("signing access token: %w", err)
	}

	return accessToken, nil
}

func (s *Service) RevokeToken(ctx context.Context, refreshTokenStr string) error {
	return s.refreshTokenRepo.Delete(ctx, refreshTokenStr)
}

func (s *Service) createRefreshToken(userID string) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	rt := &auth.RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := s.refreshTokenRepo.Create(context.Background(), rt); err != nil {
		return "", fmt.Errorf("storing refresh token: %w", err)
	}

	return token, nil
}
