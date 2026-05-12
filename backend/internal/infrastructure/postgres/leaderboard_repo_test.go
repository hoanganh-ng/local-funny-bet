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
	db := setupTestDB(t)

	ctx := context.Background()
	lbRepo := NewLeaderboardRepo(db)
	predRepo := NewPredictionRepo(db)

	userAID := uuid.New().String()
	userBID := uuid.New().String()
	userCID := uuid.New().String()
	leaderboardID := uuid.New().String()

	seedTestUser(t, db, userAID)
	seedTestUser(t, db, userBID)
	seedTestUser(t, db, userCID)

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

	memberC := &leaderboard.Member{
		ID:            uuid.New().String(),
		LeaderboardID: leaderboardID,
		UserID:        userCID,
		Role:          leaderboard.RoleMember,
		JoinedAt:      time.Now(),
	}
	if err := lbRepo.AddMember(ctx, memberC); err != nil {
		t.Fatalf("failed to add member C: %v", err)
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

	// User C predictions: all incorrect (0 points)
	predC1 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userCID,
		MatchID:   match1ID,
		Value:     prediction.ValueDraw, // match1: home won (2-1)
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predC1); err != nil {
		t.Fatalf("failed to upsert prediction C1: %v", err)
	}

	predC2 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userCID,
		MatchID:   match2ID,
		Value:     prediction.ValueHomeWin, // match2: draw (1-1)
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predC2); err != nil {
		t.Fatalf("failed to upsert prediction C2: %v", err)
	}

	predC3 := &prediction.Prediction{
		ID:        uuid.New().String(),
		UserID:    userCID,
		MatchID:   match3ID,
		Value:     prediction.ValueDraw, // match3: away won (0-3)
		UpdatedAt: time.Now(),
	}
	if err := predRepo.Upsert(ctx, predC3); err != nil {
		t.Fatalf("failed to upsert prediction C3: %v", err)
	}

	scores, err := lbRepo.GetScores(ctx, leaderboardID)
	if err != nil {
		t.Fatalf("failed to get scores: %v", err)
	}

	// Expected: A (2pts), B (1pt), C (0pts or absent)
	// Query uses COUNT(*) which excludes users with no correct predictions
	if len(scores) < 2 || len(scores) > 3 {
		t.Fatalf("expected 2 or 3 scores (A, B, optionally C), got %d", len(scores))
	}

	// User A: 2 points (first place)
	if scores[0].UserID != userAID {
		t.Errorf("expected first user to be %s, got %s", userAID, scores[0].UserID)
	}
	if scores[0].Points != 2 {
		t.Errorf("expected user A to have 2 points, got %d", scores[0].Points)
	}

	// User B: 1 point (second place)
	if scores[1].UserID != userBID {
		t.Errorf("expected second user to be %s, got %s", userBID, scores[1].UserID)
	}
	if scores[1].Points != 1 {
		t.Errorf("expected user B to have 1 point, got %d", scores[1].Points)
	}

	// User C: 0 points (absent from results, since query uses COUNT with WHERE on correct predictions)
	// GetScores uses COUNT(*) which only returns users with at least one correct prediction
	if len(scores) == 3 {
		t.Errorf("user C with 0 points should be absent from results, but got %d users", len(scores))
	}

	_, err = db.Exec("DELETE FROM predictions WHERE user_id IN ($1, $2, $3)", userAID, userBID, userCID)
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
	cleanupTestUser(t, db, userCID)
}
