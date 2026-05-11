package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"wc2026/internal/domain/auth"
	"wc2026/internal/domain/user"
)

type Service struct {
	userRepo       user.Repository
	provider       auth.Provider
	tokenSigner    auth.TokenSigner
	refreshSecret  string
	refreshTokens  map[string]refreshToken
	refreshTokenMu sync.RWMutex
}

type refreshToken struct {
	userID    string
	expiresAt time.Time
}

func NewService(userRepo user.Repository, provider auth.Provider, tokenSigner auth.TokenSigner, refreshSecret string) *Service {
	return &Service{
		userRepo:      userRepo,
		provider:      provider,
		tokenSigner:   tokenSigner,
		refreshSecret: refreshSecret,
		refreshTokens: make(map[string]refreshToken),
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
	s.refreshTokenMu.RLock()
	rt, ok := s.refreshTokens[refreshTokenStr]
	s.refreshTokenMu.RUnlock()

	if !ok {
		return "", fmt.Errorf("invalid refresh token: %w", fmt.Errorf("unauthorized"))
	}

	if time.Now().After(rt.expiresAt) {
		s.refreshTokenMu.Lock()
		delete(s.refreshTokens, refreshTokenStr)
		s.refreshTokenMu.Unlock()
		return "", fmt.Errorf("refresh token expired: %w", fmt.Errorf("unauthorized"))
	}

	u, err := s.userRepo.GetByID(ctx, rt.userID)
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
	s.refreshTokenMu.Lock()
	delete(s.refreshTokens, refreshTokenStr)
	s.refreshTokenMu.Unlock()
	return nil
}

func (s *Service) createRefreshToken(userID string) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	s.refreshTokenMu.Lock()
	s.refreshTokens[token] = refreshToken{
		userID:    userID,
		expiresAt: expiresAt,
	}
	s.refreshTokenMu.Unlock()

	return token, nil
}
