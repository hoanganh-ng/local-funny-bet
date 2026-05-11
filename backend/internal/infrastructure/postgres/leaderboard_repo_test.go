package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"wc2026/internal/domain/leaderboard"
	"wc2026/internal/domain/prediction"
)

func TestLeaderboardRepo_GetScores(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()
	lbRepo := NewLeaderboardRepo(db)
	predRepo := NewPredictionRepo(db)

	userAID := uuid.New().String()
	userBID := uuid.New().String()
	leaderboardID := uuid.New().String()

	seedTestUser(t, db, userAID)
	seedTestUser(t, db, userBID)

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

	match1ID := uuid.New().String()
	match2ID := uuid.New().String()
	match3ID := uuid.New().String()

	homeScore1 := 2
	awayScore1 := 1
	homeScore2 := 1
	awayScore2 := 1
	homeScore3 := 0
	awayScore3 := 3

	_, err = db.Exec(`INSERT INTO matches (id, tournament_id, home_team, away_team, home_score, away_score, kickoff_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		match1ID, tournamentID, "Team A", "Team B", homeScore1, awayScore1, time.Now().Add(-2*time.Hour), "finished")
	if err != nil {
		t.Fatalf("failed to seed match 1: %v", err)
	}

	_, err = db.Exec(`INSERT INTO matches (id, tournament_id, home_team, away_team, home_score, away_score, kickoff_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		match2ID, tournamentID, "Team C", "Team D", homeScore2, awayScore2, time.Now().Add(-2*time.Hour), "finished")
	if err != nil {
		t.Fatalf("failed to seed match 2: %v", err)
	}

	_, err = db.Exec(`INSERT INTO matches (id, tournament_id, home_team, away_team, home_score, away_score, kickoff_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		match3ID, tournamentID, "Team E", "Team F", homeScore3, awayScore3, time.Now().Add(-2*time.Hour), "finished")
	if err != nil {
		t.Fatalf("failed to seed match 3: %v", err)
	}

	lb := &leaderboard.Leaderboard{
		ID:        leaderboardID,
		Name:      "Test Leaderboard",
		CreatedBy: userAID,
		CreatedAt: time.Now(),
	}
	if err := lbRepo.Create(ctx, lb); err != nil {
		t.Fatalf("failed to create leaderboard: %v", err)
	}

	memberA := &leaderboard.Member{
		ID:            uuid.New().String(),
		LeaderboardID: leaderboardID,
		UserID:        userAID,
		Role:          leaderboard.RoleOwner,
		JoinedAt:      time.Now(),
	}
	if err := lbRepo.AddMember(ctx, memberA); err != nil {
		t.Fatalf("failed to add member A: %v", err)
	}

	memberB := &leaderboard.Member{
		ID:            uuid.New().String(),
		LeaderboardID: leaderboardID,
		UserID:        userBID,
		Role:          leaderboard.RoleMember,
		JoinedAt:      time.Now(),
	}
	if err := lbRepo.AddMember(ctx, memberB); err != nil {
		t.Fatalf("failed to add member B: %v", err)
	}

	predA1 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userAID,
		MatchID:   match1ID,
		Value:     prediction.ValueHomeWin,
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predA1); err != nil {
		t.Fatalf("failed to upsert prediction A1: %v", err)
	}

	predA2 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userAID,
		MatchID:   match2ID,
		Value:     prediction.ValueDraw,
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predA2); err != nil {
		t.Fatalf("failed to upsert prediction A2: %v", err)
	}

	predA3 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userAID,
		MatchID:   match3ID,
		Value:     prediction.ValueHomeWin,
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predA3); err != nil {
		t.Fatalf("failed to upsert prediction A3: %v", err)
	}

	predB1 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userBID,
		MatchID:   match1ID,
		Value:     prediction.ValueAwayWin,
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predB1); err != nil {
		t.Fatalf("failed to upsert prediction B1: %v", err)
	}

	predB2 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userBID,
		MatchID:   match2ID,
		Value:     prediction.ValueDraw,
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predB2); err != nil {
		t.Fatalf("failed to upsert prediction B2: %v", err)
	}

	predB3 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userBID,
		MatchID:   match3ID,
		Value:     prediction.ValueHomeWin,
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predB3); err != nil {
		t.Fatalf("failed to upsert prediction B3: %v", err)
	}

	scores, err := lbRepo.GetScores(ctx, leaderboardID)
	if err != nil {
		t.Fatalf("failed to get scores: %v", err)
	}

	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}

	if scores[0].UserID != userAID {
		t.Errorf("expected first user to be %s, got %s", userAID, scores[0].UserID)
	}
	if scores[0].Points != 2 {
		t.Errorf("expected user A to have 2 points, got %d", scores[0].Points)
	}

	if scores[1].UserID != userBID {
		t.Errorf("expected second user to be %s, got %s", userBID, scores[1].UserID)
	}
	if scores[1].Points != 1 {
		t.Errorf("expected user B to have 1 point, got %d", scores[1].Points)
	}

	_, err = db.Exec("DELETE FROM predictions WHERE user_id IN ($1, $2)", userAID, userBID)
	if err != nil {
		t.Logf("failed to cleanup predictions: %v", err)
	}

	_, err = db.Exec("DELETE FROM leaderboard_members WHERE leaderboard_id = $1", leaderboardID)
	if err != nil {
		t.Logf("failed to cleanup leaderboard members: %v", err)
	}

	_, err = db.Exec("DELETE FROM leaderboards WHERE id = $1", leaderboardID)
	if err != nil {
		t.Logf("failed to cleanup leaderboard: %v", err)
	}

	_, err = db.Exec("DELETE FROM matches WHERE id IN ($1, $2, $3)", match1ID, match2ID, match3ID)
	if err != nil {
		t.Logf("failed to cleanup matches: %v", err)
	}

	cleanupTestUser(t, db, userAID)
	cleanupTestUser(t, db, userBID)
}
