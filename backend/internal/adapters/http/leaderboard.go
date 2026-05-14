package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"wc2026/internal/domain"
	"wc2026/internal/domain/auth"
	"wc2026/internal/domain/leaderboard"
	"wc2026/internal/middleware"
)

type LeaderboardService interface {
	CreateLeaderboard(ctx context.Context, userID, name string) (*leaderboard.Leaderboard, error)
	GetLeaderboard(ctx context.Context, id string) (*leaderboard.Leaderboard, error)
	GetScores(ctx context.Context, id string) ([]*leaderboard.Score, error)
	ListMyLeaderboards(ctx context.Context, userID string) ([]*leaderboard.Leaderboard, error)
	GenerateInvite(ctx context.Context, leaderboardID, userID string) (token string, expiresAt time.Time, err error)
	JoinLeaderboard(ctx context.Context, userID, inviteToken string) (string, error)
}

type LeaderboardHandler struct {
	service LeaderboardService
}

func NewLeaderboardHandler(service LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

func (h *LeaderboardHandler) CreateLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}

	lb, err := h.service.CreateLeaderboard(r.Context(), userID, req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, lb)
}

func (h *LeaderboardHandler) ListLeaderboards(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	leaderboards, err := h.service.ListMyLeaderboards(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, leaderboards)
}

func (h *LeaderboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "leaderboard ID is required")
		return
	}

	lb, err := h.service.GetLeaderboard(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "leaderboard not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, lb)
}

func (h *LeaderboardHandler) GetScores(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "leaderboard ID is required")
		return
	}

	scores, err := h.service.GetScores(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "leaderboard not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if scores == nil {
		scores = []*leaderboard.Score{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"scores":       scores,
		"member_count": len(scores),
	})
}

func (h *LeaderboardHandler) GenerateInvite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	leaderboardID := r.PathValue("id")
	if leaderboardID == "" {
		respondError(w, http.StatusBadRequest, "leaderboard ID is required")
		return
	}

	token, expiresAt, err := h.service.GenerateInvite(r.Context(), leaderboardID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "only owner can generate invites")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"token":      token,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

func (h *LeaderboardHandler) JoinLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		InviteToken string `json:"invite_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.InviteToken == "" {
		respondError(w, http.StatusBadRequest, "invite_token is required")
		return
	}

	leaderboardID, err := h.service.JoinLeaderboard(r.Context(), userID, req.InviteToken)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrExpired):
			respondError(w, http.StatusUnauthorized, "invite token expired")
		case errors.Is(err, auth.ErrInvalid):
			respondError(w, http.StatusBadRequest, "invite token invalid")
		case errors.Is(err, domain.ErrAlreadyExists):
			// leaderboardID is set even on ErrAlreadyExists; fall back to token parse if not
			if leaderboardID == "" {
				parts := strings.Split(req.InviteToken, ".")
				if len(parts) >= 1 {
					leaderboardID = parts[0]
				}
			}
			respondJSON(w, http.StatusConflict, map[string]string{
				"error":          "already a member",
				"leaderboard_id": leaderboardID,
			})
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"leaderboard_id": leaderboardID})
}
