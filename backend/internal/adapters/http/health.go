package http

import (
	"database/sql"
	"net/http"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(); err != nil {
		respondError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
