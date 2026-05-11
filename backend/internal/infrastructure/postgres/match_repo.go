package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wc2026/internal/domain"
	"wc2026/internal/domain/match"
)

type MatchRepo struct {
	db *sql.DB
}

func NewMatchRepo(db *sql.DB) *MatchRepo {
	return &MatchRepo{db: db}
}

func (r *MatchRepo) GetByID(ctx context.Context, id string) (*match.Match, error) {
	query := `
		SELECT id, tournament_id, home_team, away_team, home_score, away_score,
		       kickoff_at, status, external_id
		FROM matches
		WHERE id = $1
	`

	var m match.Match
	var homeScore, awayScore sql.NullInt64
	var externalID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.TournamentID,
		&m.HomeTeam,
		&m.AwayTeam,
		&homeScore,
		&awayScore,
		&m.KickoffAt,
		&m.Status,
		&externalID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting match by id %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting match by id %s: %w", id, err)
	}

	if homeScore.Valid {
		score := int(homeScore.Int64)
		m.HomeScore = &score
	}
	if awayScore.Valid {
		score := int(awayScore.Int64)
		m.AwayScore = &score
	}
	if externalID.Valid {
		m.ExternalID = &externalID.String
	}

	return &m, nil
}

func (r *MatchRepo) List(ctx context.Context, status string) ([]*match.Match, error) {
	query := `
		SELECT id, tournament_id, home_team, away_team, home_score, away_score,
		       kickoff_at, status, external_id
		FROM matches
	`

	args := []interface{}{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}

	query += " ORDER BY kickoff_at ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing matches: %w", err)
	}
	defer rows.Close()

	var matches []*match.Match
	for rows.Next() {
		var m match.Match
		var homeScore, awayScore sql.NullInt64
		var externalID sql.NullString

		err := rows.Scan(
			&m.ID,
			&m.TournamentID,
			&m.HomeTeam,
			&m.AwayTeam,
			&homeScore,
			&awayScore,
			&m.KickoffAt,
			&m.Status,
			&externalID,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning match: %w", err)
		}

		if homeScore.Valid {
			score := int(homeScore.Int64)
			m.HomeScore = &score
		}
		if awayScore.Valid {
			score := int(awayScore.Int64)
			m.AwayScore = &score
		}
		if externalID.Valid {
			m.ExternalID = &externalID.String
		}

		matches = append(matches, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating matches: %w", err)
	}

	return matches, nil
}

func (r *MatchRepo) Upsert(ctx context.Context, m *match.Match) error {
	query := `
		INSERT INTO matches (id, tournament_id, home_team, away_team, home_score, away_score,
		                     kickoff_at, status, external_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			tournament_id = EXCLUDED.tournament_id,
			home_team = EXCLUDED.home_team,
			away_team = EXCLUDED.away_team,
			home_score = EXCLUDED.home_score,
			away_score = EXCLUDED.away_score,
			kickoff_at = EXCLUDED.kickoff_at,
			status = EXCLUDED.status,
			external_id = EXCLUDED.external_id
	`

	var homeScore, awayScore sql.NullInt64
	if m.HomeScore != nil {
		homeScore = sql.NullInt64{Int64: int64(*m.HomeScore), Valid: true}
	}
	if m.AwayScore != nil {
		awayScore = sql.NullInt64{Int64: int64(*m.AwayScore), Valid: true}
	}

	var externalID sql.NullString
	if m.ExternalID != nil {
		externalID = sql.NullString{String: *m.ExternalID, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		m.ID,
		m.TournamentID,
		m.HomeTeam,
		m.AwayTeam,
		homeScore,
		awayScore,
		m.KickoffAt,
		m.Status,
		externalID,
	)

	if err != nil {
		return fmt.Errorf("upserting match: %w", err)
	}

	return nil
}
