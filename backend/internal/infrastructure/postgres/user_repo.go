package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wc2026/internal/domain"
	"wc2026/internal/domain/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, email, name, avatar_url, created_at
		FROM users
		WHERE id = $1
	`

	var u user.User
	var avatarURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&avatarURL,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting user by id %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting user by id %s: %w", id, err)
	}

	if avatarURL.Valid {
		u.AvatarURL = &avatarURL.String
	}

	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, name, avatar_url, created_at
		FROM users
		WHERE email = $1
	`

	var u user.User
	var avatarURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&avatarURL,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting user by email %s: %w", email, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting user by email %s: %w", email, err)
	}

	if avatarURL.Valid {
		u.AvatarURL = &avatarURL.String
	}

	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, email, name, avatar_url, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	var avatarURL sql.NullString
	if u.AvatarURL != nil {
		avatarURL = sql.NullString{String: *u.AvatarURL, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		u.ID,
		u.Email,
		u.Name,
		avatarURL,
		u.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}
