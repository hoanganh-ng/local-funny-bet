package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrExpired = errors.New("invite token expired")
	ErrInvalid = errors.New("invite token invalid")
)

// SignInviteToken creates a stateless HMAC-signed invite token.
// Format: leaderboardID.expiresAt.signature
func SignInviteToken(leaderboardID string, secret string) string {
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()
	payload := fmt.Sprintf("%s.%d", leaderboardID, expiresAt)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s.%s", payload, signature)
}

// VerifyInviteToken validates the signature and expiration.
// Returns leaderboardID if valid.
func VerifyInviteToken(token string, secret string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalid
	}

	leaderboardID := parts[0]
	expiresAtStr := parts[1]
	providedSignature := parts[2]

	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return "", ErrInvalid
	}

	if time.Now().Unix() > expiresAt {
		return "", ErrExpired
	}

	payload := fmt.Sprintf("%s.%d", leaderboardID, expiresAt)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(providedSignature)) {
		return "", ErrInvalid
	}

	return leaderboardID, nil
}
