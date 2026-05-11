package auth

import "time"

type TokenSigner interface {
	Sign(userID, email, name, provider string, ttl time.Duration) (string, error)
}

type TokenVerifier interface {
	Verify(token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID   string
	Email    string
	Name     string
	Provider string
}
