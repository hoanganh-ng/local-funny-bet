package http

import (
	"context"
	"errors"
	"net/http"

	"wc2026/internal/domain"
	"wc2026/internal/domain/match"
)

type MatchService interface {
	ListMatches(ctx context.Context, status string) ([]*match.Match, error)
	GetMatch(ctx context.Context, id string) (*match.Match, error)
}

type MatchHandler struct {
	service MatchService
}

func NewMatchHandler(service MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) ListMatches(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	matches, err := h.service.ListMatches(r.Context(), status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, matches)
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
