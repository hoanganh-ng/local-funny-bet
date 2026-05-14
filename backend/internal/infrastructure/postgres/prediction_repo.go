package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wc2026/internal/domain"
	"wc2026/internal/domain/prediction"
)

type PredictionRepo struct {
	db *sql.DB
}

func NewPredictionRepo(db *sql.DB) *PredictionRepo {
	return &PredictionRepo{db: db}
}

func (r *PredictionRepo) Upsert(ctx context.Context, p *prediction.Prediction) error {
	query := `
		INSERT INTO predictions (id, user_id, match_id, value, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, match_id) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		p.ID,
		p.UserID,
		p.MatchID,
		p.Value,
		p.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("upserting prediction: %w", err)
	}

	return nil
}

func (r *PredictionRepo) GetByUserAndMatch(ctx context.Context, userID, matchID string) (*prediction.Prediction, error) {
	query := `
		SELECT id, user_id, match_id, value, updated_at
		FROM predictions
		WHERE user_id = $1 AND match_id = $2
	`

	var p prediction.Prediction
	err := r.db.QueryRowContext(ctx, query, userID, matchID).Scan(
		&p.ID,
		&p.UserID,
		&p.MatchID,
		&p.Value,
		&p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting prediction for user %s match %s: %w", userID, matchID, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting prediction for user %s match %s: %w", userID, matchID, err)
	}

	return &p, nil
}

func (r *PredictionRepo) ListByMatch(ctx context.Context, matchID string) ([]*prediction.Prediction, error) {
	query := `
		SELECT p.id, p.user_id, p.match_id, p.value, p.updated_at,
		       COALESCE(u.name, '') AS user_name, u.avatar_url
		FROM predictions p
		LEFT JOIN users u ON u.id = p.user_id
		WHERE p.match_id = $1
		ORDER BY p.updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, matchID)
	if err != nil {
		return nil, fmt.Errorf("listing predictions for match %s: %w", matchID, err)
	}
	defer rows.Close()

	var predictions []*prediction.Prediction
	for rows.Next() {
		var p prediction.Prediction
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.MatchID,
			&p.Value,
			&p.UpdatedAt,
			&p.UserName,
			&p.UserAvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning prediction: %w", err)
		}

		predictions = append(predictions, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating predictions: %w", err)
	}

	return predictions, nil
}

func (r *PredictionRepo) ListByUser(ctx context.Context, userID string) ([]*prediction.PredictionWithMatch, error) {
	query := `
		SELECT p.id, p.user_id, p.match_id, p.value, p.updated_at,
		       m.home_team, m.away_team, m.kickoff_at, m.status,
		       m.home_score, m.away_score
		FROM predictions p
		JOIN matches m ON m.id = p.match_id
		WHERE p.user_id = $1
		ORDER BY m.kickoff_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("listing predictions for user %s: %w", userID, err)
	}
	defer rows.Close()

	var results []*prediction.PredictionWithMatch
	for rows.Next() {
		var p prediction.PredictionWithMatch
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.MatchID,
			&p.Value,
			&p.UpdatedAt,
			&p.HomeTeam,
			&p.AwayTeam,
			&p.KickoffAt,
			&p.Status,
			&p.HomeScore,
			&p.AwayScore,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning prediction history: %w", err)
		}

		results = append(results, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating prediction history: %w", err)
	}

	return results, nil
}
