package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"wc2026/internal/domain"
	"wc2026/internal/domain/prediction"
	"wc2026/internal/middleware"
)

type PredictionService interface {
	UpsertPrediction(ctx context.Context, userID, matchID, value string) error
	ListByMatch(ctx context.Context, leaderboardID, matchID, userID string) ([]*prediction.Prediction, error)
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

	respondJSON(w, http.StatusOK, predictions)
}
