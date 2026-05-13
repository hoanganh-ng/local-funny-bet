package http

import (
	"context"
	"net/http"

	"wc2026/internal/domain/tournament"
)

type TournamentRepository interface {
	GetActive(ctx context.Context) (*tournament.Tournament, error)
}

type TournamentHandler struct {
	repo TournamentRepository
}

func NewTournamentHandler(repo TournamentRepository) *TournamentHandler {
	return &TournamentHandler{repo: repo}
}

func (h *TournamentHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	t, err := h.repo.GetActive(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, t)
}
