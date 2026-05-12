package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"wc2026/internal/domain"
	"wc2026/internal/middleware"
)

type AuthService interface {
	GetAuthURL(state string) string
	HandleGoogleCallback(ctx context.Context, code string) (accessToken, refreshToken string, err error)
	GetUserByID(ctx context.Context, userID string) (id, email, name string, avatarURL *string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	RevokeToken(ctx context.Context, refreshToken string) error
}

type AuthHandler struct {
	service     AuthService
	frontendURL string
}

func NewAuthHandler(service AuthService, frontendURL string) *AuthHandler {
	return &AuthHandler{
		service:     service,
		frontendURL: frontendURL,
	}
}

func (h *AuthHandler) InitiateGoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := h.service.GetAuthURL("random-state")
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, h.frontendURL+"/login?error=missing_code", http.StatusTemporaryRedirect)
		return
	}

	_, refreshToken, err := h.service.HandleGoogleCallback(r.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			http.Redirect(w, r, h.frontendURL+"/login?error=unauthorized", http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, h.frontendURL+"/login?error=server_error", http.StatusTemporaryRedirect)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to frontend OAuth callback handler
	redirectURL := fmt.Sprintf("%s/auth/callback?success=true", h.frontendURL)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) ExchangeRefreshForAccess(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		respondError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	accessToken, err := h.service.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		respondError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	accessToken, err := h.service.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, email, name, avatarURL, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":         id,
		"email":      email,
		"name":       name,
		"avatar_url": avatarURL,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		if err := h.service.RevokeToken(r.Context(), cookie.Value); err != nil {
			fmt.Printf("error revoking token: %v\n", err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
