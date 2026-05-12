package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"wc2026/internal/domain"
	"wc2026/internal/domain/auth"
)

type RefreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, rt *auth.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (token, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		rt.Token,
		rt.UserID,
		rt.ExpiresAt,
		rt.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("creating refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepo) GetByToken(ctx context.Context, token string) (*auth.RefreshToken, error) {
	query := `
		SELECT token, user_id, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`

	var rt auth.RefreshToken

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&rt.Token,
		&rt.UserID,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting refresh token: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting refresh token: %w", err)
	}

	return &rt, nil
}

func (r *RefreshTokenRepo) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`

	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("deleting refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepo) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("deleting expired refresh tokens: %w", err)
	}

	return nil
}
