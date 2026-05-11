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
		SELECT id, user_id, match_id, value, updated_at
		FROM predictions
		WHERE match_id = $1
		ORDER BY updated_at DESC
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
