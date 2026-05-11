package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSignAndVerifyInviteToken(t *testing.T) {
	secret := "test-secret-key"
	leaderboardID := "lb-12345"

	token := SignInviteToken(leaderboardID, secret)

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("token should have 3 parts, got %d", len(parts))
	}

	gotID, err := VerifyInviteToken(token, secret)
	if err != nil {
		t.Fatalf("VerifyInviteToken() error = %v", err)
	}
	if gotID != leaderboardID {
		t.Errorf("VerifyInviteToken() ID = %v, want %v", gotID, leaderboardID)
	}
}

func TestVerifyInviteToken_Expiration(t *testing.T) {
	secret := "test-secret-key"

	expiredAt := time.Now().Add(-1 * time.Hour).Unix()
	payload := fmt.Sprintf("lb-12345.%d", expiredAt)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	expiredToken := payload + "." + signature

	_, err := VerifyInviteToken(expiredToken, secret)
	if err != ErrExpired {
		t.Errorf("VerifyInviteToken() error = %v, want ErrExpired", err)
	}
}

func TestVerifyInviteToken_Tampering(t *testing.T) {
	secret := "test-secret-key"
	leaderboardID := "lb-12345"

	token := SignInviteToken(leaderboardID, secret)

	parts := strings.Split(token, ".")
	parts[0] = "lb-99999"
	tamperedToken := strings.Join(parts, ".")

	_, err := VerifyInviteToken(tamperedToken, secret)
	if err != ErrInvalid {
		t.Errorf("VerifyInviteToken() error = %v, want ErrInvalid", err)
	}
}

func TestVerifyInviteToken_InvalidFormat(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"missing parts", "onlyonepart"},
		{"empty token", ""},
		{"two parts only", "part1.part2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyInviteToken(tt.token, "secret")
			if err != ErrInvalid {
				t.Errorf("VerifyInviteToken() error = %v, want ErrInvalid", err)
			}
		})
	}
}
