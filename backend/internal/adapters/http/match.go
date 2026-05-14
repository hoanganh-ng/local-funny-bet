package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"wc2026/internal/domain"
	"wc2026/internal/domain/match"
)

type MatchService interface {
	ListMatches(ctx context.Context, status string, limit int, after time.Time) ([]*match.Match, error)
	GetMatch(ctx context.Context, id string) (*match.Match, error)
}

type matchListResponse struct {
	Matches    []*match.Match `json:"matches"`
	NextCursor *string        `json:"next_cursor,omitempty"`
	HasMore    bool           `json:"has_more"`
}

type MatchHandler struct {
	service MatchService
}

func NewMatchHandler(service MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) ListMatches(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 {
			respondError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		if parsed > 50 {
			parsed = 50
		}
		limit = parsed
	}

	var after time.Time
	if a := r.URL.Query().Get("after"); a != "" {
		parsed, err := time.Parse(time.RFC3339, a)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid after timestamp")
			return
		}
		after = parsed
	}

	matches, err := h.service.ListMatches(r.Context(), status, limit, after)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := matchListResponse{
		Matches: matches,
		HasMore: len(matches) == limit,
	}
	if resp.HasMore && len(matches) > 0 {
		cursor := matches[len(matches)-1].KickoffAt.UTC().Format(time.RFC3339)
		resp.NextCursor = &cursor
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *MatchHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "missing match id")
		return
	}

	m, err := h.service.GetMatch(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "match not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, m)
}
