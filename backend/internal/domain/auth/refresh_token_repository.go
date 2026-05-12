package auth

import (
	"context"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *RefreshToken) error
	GetByToken(ctx context.Context, token string) (*RefreshToken, error)
	Delete(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}
