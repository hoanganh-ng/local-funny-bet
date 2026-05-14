package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"wc2026/internal/domain/auth"
	"wc2026/internal/domain/leaderboard"
	"wc2026/internal/middleware"
)

// mockTokenVerifier for testing JWT middleware
type mockTokenVerifier struct {
	verifyFunc func(token string) (*auth.TokenClaims, error)
}

func (m *mockTokenVerifier) Verify(token string) (*auth.TokenClaims, error) {
	return m.verifyFunc(token)
}

func TestJWTMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		verifyFunc     func(token string) (*auth.TokenClaims, error)
		wantStatus     int
		wantError      string
		wantHandlerRun bool
	}{
		{
			name:       "valid JWT",
			authHeader: "Bearer valid-token",
			verifyFunc: func(token string) (*auth.TokenClaims, error) {
				return &auth.TokenClaims{UserID: "user-123"}, nil
			},
			wantStatus:     http.StatusOK,
			wantHandlerRun: true,
		},
		{
			name:       "expired JWT",
			authHeader: "Bearer expired-token",
			verifyFunc: func(token string) (*auth.TokenClaims, error) {
				return nil, auth.ErrExpired
			},
			wantStatus:     http.StatusUnauthorized,
			wantError:      "unauthorized",
			wantHandlerRun: false,
		},
		{
			name:       "missing Authorization header",
			authHeader: "",
			verifyFunc: func(token string) (*auth.TokenClaims, error) {
				return nil, nil
			},
			wantStatus:     http.StatusUnauthorized,
			wantError:      "unauthorized",
			wantHandlerRun: false,
		},
		{
			name:       "malformed token - not valid format",
			authHeader: "Bearer not-valid-base64-!!!",
			verifyFunc: func(token string) (*auth.TokenClaims, error) {
				return nil, auth.ErrInvalid
			},
			wantStatus:     http.StatusUnauthorized,
			wantError:      "unauthorized",
			wantHandlerRun: false,
		},
		{
			name:       "malformed token - missing Bearer prefix",
			authHeader: "token-without-bearer",
			verifyFunc: func(token string) (*auth.TokenClaims, error) {
				return nil, nil
			},
			wantStatus:     http.StatusUnauthorized,
			wantError:      "unauthorized",
			wantHandlerRun: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerRun := false
			var capturedUserID string

			// Test handler that captures context values
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerRun = true
				userID, ok := r.Context().Value(middleware.UserIDKey).(string)
				if ok {
					capturedUserID = userID
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			})

			// Create mock verifier
			verifier := &mockTokenVerifier{verifyFunc: tt.verifyFunc}

			// Wrap handler with middleware
			wrappedHandler := middleware.RequireAuth(verifier)(testHandler)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Record response
			w := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(w, req)

			// Assert status code
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			// Assert handler execution
			if handlerRun != tt.wantHandlerRun {
				t.Errorf("handlerRun = %v, want %v", handlerRun, tt.wantHandlerRun)
			}

			// Assert error response body
			if tt.wantError != "" {
				var resp map[string]string
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp["error"] != tt.wantError {
					t.Errorf("error = %q, want %q", resp["error"], tt.wantError)
				}
			}

			// Assert userID in context for valid token
			if tt.wantHandlerRun && capturedUserID != "user-123" {
				t.Errorf("capturedUserID = %q, want %q", capturedUserID, "user-123")
			}
		})
	}
}

func TestJoinLeaderboard_RejectsQueryStringToken(t *testing.T) {
	// Mock service
	svc := &mockLeaderboardService{
		joinFunc: func(ctx context.Context, userID, inviteToken string) (string, error) {
			t.Fatal("JoinLeaderboard should not be called when token is in query string")
			return "", nil
		},
	}

	handler := NewLeaderboardHandler(svc)

	// Create request with token in query string
	req := httptest.NewRequest(http.MethodPost, "/leaderboards/join?invite_token=xxx", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")

	// Inject userID into context (simulate middleware)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.JoinLeaderboard(w, req)

	// Assert 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "invite_token is required" {
		t.Errorf("error = %q, want %q", resp["error"], "invite_token is required")
	}
}

// mockLeaderboardService for testing
type mockLeaderboardService struct {
	joinFunc func(ctx context.Context, userID, inviteToken string) (string, error)
}

func (m *mockLeaderboardService) CreateLeaderboard(ctx context.Context, userID, name string) (*leaderboard.Leaderboard, error) {
	return nil, nil
}

func (m *mockLeaderboardService) GetLeaderboard(ctx context.Context, id string) (*leaderboard.Leaderboard, error) {
	return nil, nil
}

func (m *mockLeaderboardService) ListMyLeaderboards(ctx context.Context, userID string) ([]*leaderboard.Leaderboard, error) {
	return nil, nil
}

func (m *mockLeaderboardService) GenerateInvite(ctx context.Context, leaderboardID, userID string) (token string, expiresAt time.Time, err error) {
	return "", time.Time{}, nil
}

func (m *mockLeaderboardService) JoinLeaderboard(ctx context.Context, userID, inviteToken string) (string, error) {
	return m.joinFunc(ctx, userID, inviteToken)
}
