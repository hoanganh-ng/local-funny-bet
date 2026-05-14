package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"wc2026/internal/domain"
	"wc2026/internal/domain/leaderboard"
)

type LeaderboardRepo struct {
	db *sql.DB
}

func NewLeaderboardRepo(db *sql.DB) *LeaderboardRepo {
	return &LeaderboardRepo{db: db}
}

func (r *LeaderboardRepo) Create(ctx context.Context, lb *leaderboard.Leaderboard) error {
	query := `
		INSERT INTO leaderboards (id, name, created_by, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		lb.ID,
		lb.Name,
		lb.CreatedBy,
		lb.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("creating leaderboard: %w", err)
	}

	return nil
}

func (r *LeaderboardRepo) GetByID(ctx context.Context, id string) (*leaderboard.Leaderboard, error) {
	query := `
		SELECT id, name, created_by, created_at
		FROM leaderboards
		WHERE id = $1
	`

	var lb leaderboard.Leaderboard
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lb.ID,
		&lb.Name,
		&lb.CreatedBy,
		&lb.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting leaderboard by id %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting leaderboard by id %s: %w", id, err)
	}

	return &lb, nil
}

func (r *LeaderboardRepo) ListByUser(ctx context.Context, userID string) ([]*leaderboard.Leaderboard, error) {
	query := `
		SELECT l.id, l.name, l.created_by, l.created_at
		FROM leaderboards l
		JOIN leaderboard_members lm ON l.id = lm.leaderboard_id
		WHERE lm.user_id = $1
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("listing leaderboards for user %s: %w", userID, err)
	}
	defer rows.Close()

	var leaderboards []*leaderboard.Leaderboard
	for rows.Next() {
		var lb leaderboard.Leaderboard
		err := rows.Scan(
			&lb.ID,
			&lb.Name,
			&lb.CreatedBy,
			&lb.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning leaderboard: %w", err)
		}

		leaderboards = append(leaderboards, &lb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating leaderboards: %w", err)
	}

	return leaderboards, nil
}

func (r *LeaderboardRepo) AddMember(ctx context.Context, m *leaderboard.Member) error {
	query := `
		INSERT INTO leaderboard_members (id, leaderboard_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		m.ID,
		m.LeaderboardID,
		m.UserID,
		m.Role,
		m.JoinedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("adding member to leaderboard: %w", err)
	}

	return nil
}

func (r *LeaderboardRepo) GetMember(ctx context.Context, leaderboardID, userID string) (*leaderboard.Member, error) {
	query := `
		SELECT id, leaderboard_id, user_id, role, joined_at
		FROM leaderboard_members
		WHERE leaderboard_id = $1 AND user_id = $2
	`

	var m leaderboard.Member
	err := r.db.QueryRowContext(ctx, query, leaderboardID, userID).Scan(
		&m.ID,
		&m.LeaderboardID,
		&m.UserID,
		&m.Role,
		&m.JoinedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting member for leaderboard %s user %s: %w", leaderboardID, userID, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("getting member for leaderboard %s user %s: %w", leaderboardID, userID, err)
	}

	return &m, nil
}

func (r *LeaderboardRepo) GetScores(ctx context.Context, leaderboardID string) ([]*leaderboard.Score, error) {
	// LEFT JOIN so members with 0 correct predictions still appear in the standings.
	query := `
		SELECT u.id, u.name, COUNT(CASE WHEN
		    m.status = 'finished' AND (
		        (m.home_score > m.away_score AND p.value = 'home') OR
		        (m.home_score = m.away_score AND p.value = 'draw') OR
		        (m.home_score < m.away_score AND p.value = 'away')
		    ) THEN 1 END) AS points
		FROM leaderboard_members lm
		JOIN users u ON u.id = lm.user_id
		LEFT JOIN predictions p ON p.user_id = lm.user_id
		LEFT JOIN matches m ON m.id = p.match_id
		WHERE lm.leaderboard_id = $1
		GROUP BY u.id, u.name
		ORDER BY points DESC, u.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, leaderboardID)
	if err != nil {
		return nil, fmt.Errorf("getting scores for leaderboard %s: %w", leaderboardID, err)
	}
	defer rows.Close()

	var scores []*leaderboard.Score
	for rows.Next() {
		var s leaderboard.Score
		err := rows.Scan(
			&s.UserID,
			&s.Name,
			&s.Points,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning score: %w", err)
		}

		scores = append(scores, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating scores: %w", err)
	}

	return scores, nil
}

func (r *LeaderboardRepo) GetAffectedByMatches(ctx context.Context, matchIDs []string) ([]string, error) {
	if len(matchIDs) == 0 {
		return []string{}, nil
	}

	query := `
		SELECT DISTINCT lm.leaderboard_id
		FROM leaderboard_members lm
		JOIN predictions p ON p.user_id = lm.user_id
		WHERE p.match_id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, matchIDs)
	if err != nil {
		return nil, fmt.Errorf("getting affected leaderboards: %w", err)
	}
	defer rows.Close()

	var leaderboardIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scanning leaderboard id: %w", err)
		}
		leaderboardIDs = append(leaderboardIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating leaderboard ids: %w", err)
	}

	return leaderboardIDs, nil
}
