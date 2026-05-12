package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wc2026/internal/domain"
	"wc2026/internal/domain/tournament"
)

type TournamentRepo struct {
	db *sql.DB
}

func NewTournamentRepo(db *sql.DB) *TournamentRepo {
	return &TournamentRepo{db: db}
}

func (r *TournamentRepo) GetByID(ctx context.Context, id string) (*tournament.Tournament, error) {
	query := `
		SELECT id, name, season, logo_url, status, external_id
		FROM tournaments
		WHERE id = $1
	`

	var t tournament.Tournament
	var logoURL, externalID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Name,
		&t.Season,
		&logoURL,
		&t.Status,
		&externalID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting tournament by id %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting tournament by id %s: %w", id, err)
	}

	if logoURL.Valid {
		t.LogoURL = &logoURL.String
	}
	if externalID.Valid {
		t.ExternalID = &externalID.String
	}

	return &t, nil
}

func (r *TournamentRepo) GetByExternalCode(ctx context.Context, code string) (*tournament.Tournament, error) {
	query := `
		SELECT id, name, season, logo_url, status, external_id
		FROM tournaments
		WHERE external_id = $1
		LIMIT 1
	`

	var t tournament.Tournament
	var logoURL, externalID sql.NullString

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&t.ID,
		&t.Name,
		&t.Season,
		&logoURL,
		&t.Status,
		&externalID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting tournament by external code %s: %w", code, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting tournament by external code %s: %w", code, err)
	}

	if logoURL.Valid {
		t.LogoURL = &logoURL.String
	}
	if externalID.Valid {
		t.ExternalID = &externalID.String
	}

	return &t, nil
}

func (r *TournamentRepo) List(ctx context.Context) ([]*tournament.Tournament, error) {
	query := `
		SELECT id, name, season, logo_url, status, external_id
		FROM tournaments
		ORDER BY season DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("listing tournaments: %w", err)
	}
	defer rows.Close()

	var tournaments []*tournament.Tournament
	for rows.Next() {
		var t tournament.Tournament
		var logoURL, externalID sql.NullString

		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Season,
			&logoURL,
			&t.Status,
			&externalID,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning tournament: %w", err)
		}

		if logoURL.Valid {
			t.LogoURL = &logoURL.String
		}
		if externalID.Valid {
			t.ExternalID = &externalID.String
		}

		tournaments = append(tournaments, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tournaments: %w", err)
	}

	return tournaments, nil
}

func (r *TournamentRepo) GetActive(ctx context.Context) (*tournament.Tournament, error) {
	query := `
		SELECT id, name, season, logo_url, status, external_id
		FROM tournaments
		ORDER BY season DESC
		LIMIT 1
	`

	var t tournament.Tournament
	var logoURL, externalID sql.NullString

	err := r.db.QueryRowContext(ctx, query).Scan(
		&t.ID,
		&t.Name,
		&t.Season,
		&logoURL,
		&t.Status,
		&externalID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting active tournament: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting active tournament: %w", err)
	}

	if logoURL.Valid {
		t.LogoURL = &logoURL.String
	}
	if externalID.Valid {
		t.ExternalID = &externalID.String
	}

	return &t, nil
}
