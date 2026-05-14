package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"wc2026/internal/domain"
	"wc2026/internal/domain/match"
	"wc2026/internal/domain/prediction"
	"wc2026/internal/middleware"
)

type PredictionService interface {
	UpsertPrediction(ctx context.Context, userID, matchID, value string) error
	ListByMatch(ctx context.Context, leaderboardID, matchID, userID string) ([]*prediction.Prediction, error)
	GetHistory(ctx context.Context, userID string) ([]*prediction.PredictionWithMatch, error)
}

type predictionView struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	UserName      string    `json:"user_name"`
	UserAvatarURL *string   `json:"user_avatar_url"`
	MatchID       string    `json:"match_id"`
	Outcome       string    `json:"outcome"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsCurrentUser bool      `json:"is_current_user"`
}

type historyView struct {
	ID              string    `json:"id"`
	MatchID         string    `json:"match_id"`
	HomeTeam        string    `json:"home_team"`
	AwayTeam        string    `json:"away_team"`
	KickoffAt       time.Time `json:"kickoff_at"`
	Status          string    `json:"status"`
	HomeScore       *int      `json:"home_score"`
	AwayScore       *int      `json:"away_score"`
	Outcome         string    `json:"outcome"`
	PredictionLabel string    `json:"prediction_label"`
	Result          *string   `json:"result"`
	IsCorrect       bool      `json:"is_correct"`
	Points          int       `json:"points"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func toOutcome(value string) string {
	switch value {
	case "home":
		return "home_win"
	case "away":
		return "away_win"
	default:
		return value
	}
}

func outcomeLabel(outcome string) string {
	switch outcome {
	case "home_win":
		return "Home"
	case "away_win":
		return "Away"
	case "draw":
		return "Draw"
	}
	return ""
}

func matchResultOutcome(status string, homeScore, awayScore *int) *string {
	if status != "finished" || homeScore == nil || awayScore == nil {
		return nil
	}
	var r string
	switch {
	case *homeScore > *awayScore:
		r = "home_win"
	case *homeScore < *awayScore:
		r = "away_win"
	default:
		r = "draw"
	}
	return &r
}

type PredictionHandler struct {
	service PredictionService
}

func NewPredictionHandler(service PredictionService) *PredictionHandler {
	return &PredictionHandler{service: service}
}

func (h *PredictionHandler) UpsertPrediction(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		MatchID string `json:"match_id"`
		Value   string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.MatchID == "" || req.Value == "" {
		respondError(w, http.StatusBadRequest, "match_id and value are required")
		return
	}

	err := h.service.UpsertPrediction(r.Context(), userID, req.MatchID, req.Value)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "match not found")
		case errors.Is(err, domain.ErrLocked):
			respondError(w, http.StatusLocked, "match has already started")
		default:
			respondError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "prediction saved"})
}

func (h *PredictionHandler) ListPredictions(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lbID := r.PathValue("lbID")
	matchID := r.PathValue("matchID")

	if lbID == "" || matchID == "" {
		respondError(w, http.StatusBadRequest, "leaderboard ID and match ID are required")
		return
	}

	predictions, err := h.service.ListByMatch(r.Context(), lbID, matchID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			respondError(w, http.StatusUnauthorized, "not a member of this leaderboard")
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, http.StatusNotFound, "not found")
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	views := make([]predictionView, len(predictions))
	for i, p := range predictions {
		views[i] = predictionView{
			ID:            p.ID,
			UserID:        p.UserID,
			UserName:      p.UserName,
			UserAvatarURL: p.UserAvatarURL,
			MatchID:       p.MatchID,
			Outcome:       toOutcome(p.Value),
			UpdatedAt:     p.UpdatedAt,
			IsCurrentUser: p.UserID == userID,
		}
	}

	respondJSON(w, http.StatusOK, views)
}

func (h *PredictionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	history, err := h.service.GetHistory(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	views := make([]historyView, len(history))
	for i, p := range history {
		outcome := toOutcome(p.Value)
		result := matchResultOutcome(p.Status, p.HomeScore, p.AwayScore)
		m := match.Match{Status: p.Status, HomeScore: p.HomeScore, AwayScore: p.AwayScore}
		points := prediction.CalculatePoint(m, p.Value)
		isCorrect := points == 1
		views[i] = historyView{
			ID:              p.ID,
			MatchID:         p.MatchID,
			HomeTeam:        p.HomeTeam,
			AwayTeam:        p.AwayTeam,
			KickoffAt:       p.KickoffAt,
			Status:          p.Status,
			HomeScore:       p.HomeScore,
			AwayScore:       p.AwayScore,
			Outcome:         outcome,
			PredictionLabel: outcomeLabel(outcome),
			Result:          result,
			IsCorrect:       isCorrect,
			Points:          points,
			UpdatedAt:       p.UpdatedAt,
		}
	}

	respondJSON(w, http.StatusOK, views)
}
