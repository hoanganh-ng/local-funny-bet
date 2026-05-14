package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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
		SELECT id, tournament_id, home_team, away_team, home_team_code, away_team_code,
		       home_score, away_score, kickoff_at, status, external_id
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
		&m.HomeTeamCode,
		&m.AwayTeamCode,
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

func (r *MatchRepo) List(ctx context.Context, status string, limit int, after time.Time) ([]*match.Match, error) {
	base := `
		SELECT id, tournament_id, home_team, away_team, home_team_code, away_team_code,
		       home_score, away_score, kickoff_at, status, external_id
		FROM matches
	`

	var conditions []string
	var args []interface{}

	if status != "" {
		args = append(args, status)
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}
	if !after.IsZero() {
		args = append(args, after)
		conditions = append(conditions, fmt.Sprintf("kickoff_at > $%d", len(args)))
	}

	query := base
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, limit)
	query += fmt.Sprintf(" ORDER BY kickoff_at ASC LIMIT $%d", len(args))

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
			&m.HomeTeamCode,
			&m.AwayTeamCode,
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
		INSERT INTO matches (id, tournament_id, home_team, away_team, home_team_code, away_team_code,
		                     home_score, away_score, kickoff_at, status, external_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (external_id) WHERE external_id IS NOT NULL DO UPDATE SET
			home_team = EXCLUDED.home_team,
			away_team = EXCLUDED.away_team,
			home_team_code = EXCLUDED.home_team_code,
			away_team_code = EXCLUDED.away_team_code,
			home_score = EXCLUDED.home_score,
			away_score = EXCLUDED.away_score,
			kickoff_at = EXCLUDED.kickoff_at,
			status = EXCLUDED.status
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
		m.HomeTeamCode,
		m.AwayTeamCode,
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

func (r *MatchRepo) UpsertMany(ctx context.Context, matches []*match.Match) error {
	if len(matches) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO matches (id, tournament_id, home_team, away_team, home_team_code, away_team_code,
		                     home_score, away_score, kickoff_at, status, external_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (external_id) WHERE external_id IS NOT NULL DO UPDATE SET
			home_team = EXCLUDED.home_team,
			away_team = EXCLUDED.away_team,
			home_team_code = EXCLUDED.home_team_code,
			away_team_code = EXCLUDED.away_team_code,
			home_score = EXCLUDED.home_score,
			away_score = EXCLUDED.away_score,
			kickoff_at = EXCLUDED.kickoff_at,
			status = EXCLUDED.status
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, m := range matches {
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

		_, err := stmt.ExecContext(ctx,
			m.ID,
			m.TournamentID,
			m.HomeTeam,
			m.AwayTeam,
			m.HomeTeamCode,
			m.AwayTeamCode,
			homeScore,
			awayScore,
			m.KickoffAt,
			m.Status,
			externalID,
		)
		if err != nil {
			return fmt.Errorf("upserting match %s: %w", m.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
