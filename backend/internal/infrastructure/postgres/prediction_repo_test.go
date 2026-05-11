package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"wc2026/internal/domain/prediction"
)

func getTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	db, err := NewDB(dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	return db
}

func TestPredictionUpsert_UniqueConstraint(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := NewPredictionRepo(db)

	userID := uuid.New().String()
	matchID := uuid.New().String()

	seedTestUser(t, db, userID)
	seedTestMatch(t, db, matchID)

	p1 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userID,
		MatchID:   matchID,
		Value:     prediction.ValueHomeWin,
		UpdatedAt: time.Now(),
	}

	if err := repo.Upsert(ctx, p1); err != nil {
		t.Fatalf("first upsert failed: %v", err)
	}

	p2 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userID,
		MatchID:   matchID,
		Value:     prediction.ValueAwayWin,
		UpdatedAt: time.Now(),
	}

	if err := repo.Upsert(ctx, p2); err != nil {
		t.Fatalf("second upsert failed: %v", err)
	}

	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM predictions WHERE user_id = $1 AND match_id = $2", userID, matchID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count predictions: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 prediction, got %d", count)
	}

	var value string
	err = db.QueryRowContext(ctx, "SELECT value FROM predictions WHERE user_id = $1 AND match_id = $2", userID, matchID).Scan(&value)
	if err != nil {
		t.Fatalf("failed to get prediction value: %v", err)
	}

	if value != prediction.ValueAwayWin {
		t.Errorf("expected value %s, got %s", prediction.ValueAwayWin, value)
	}

	cleanupTestPrediction(t, db, userID, matchID)
	cleanupTestMatch(t, db, matchID)
	cleanupTestUser(t, db, userID)
}

func seedTestUser(t *testing.T, db *sql.DB, userID string) {
	t.Helper()
	_, err := db.Exec("INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4)",
		userID, "test@example.com", "Test User", time.Now())
	if err != nil {
		t.Fatalf("failed to seed test user: %v", err)
	}
}

func seedTestMatch(t *testing.T, db *sql.DB, matchID string) {
	t.Helper()

	var tournamentID string
	err := db.QueryRow("SELECT id FROM tournaments LIMIT 1").Scan(&tournamentID)
	if err != nil {
		tournamentID = uuid.New().String()
		_, err = db.Exec("INSERT INTO tournaments (id, name, season, status) VALUES ($1, $2, $3, $4)",
			tournamentID, "Test Tournament", "2026", "active")
		if err != nil {
			t.Fatalf("failed to seed test tournament: %v", err)
		}
	}

	_, err = db.Exec(`INSERT INTO matches (id, tournament_id, home_team, away_team, kickoff_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		matchID, tournamentID, "Team A", "Team B", time.Now().Add(24*time.Hour), "scheduled")
	if err != nil {
		t.Fatalf("failed to seed test match: %v", err)
	}
}

func cleanupTestPrediction(t *testing.T, db *sql.DB, userID, matchID string) {
	t.Helper()
	_, err := db.Exec("DELETE FROM predictions WHERE user_id = $1 AND match_id = $2", userID, matchID)
	if err != nil {
		t.Logf("failed to cleanup test prediction: %v", err)
	}
}

func cleanupTestMatch(t *testing.T, db *sql.DB, matchID string) {
	t.Helper()
	_, err := db.Exec("DELETE FROM matches WHERE id = $1", matchID)
	if err != nil {
		t.Logf("failed to cleanup test match: %v", err)
	}
}

func cleanupTestUser(t *testing.T, db *sql.DB, userID string) {
	t.Helper()
	_, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Logf("failed to cleanup test user: %v", err)
	}
}
