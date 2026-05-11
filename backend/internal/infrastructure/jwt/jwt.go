package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"wc2026/internal/domain"
	"wc2026/internal/domain/auth"
)

type Signer struct {
	secret string
}

func NewSigner(secret string) *Signer {
	return &Signer{secret: secret}
}

type Verifier struct {
	secret string
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{secret: secret}
}

type claims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Exp      int64  `json:"exp"`
}

func (s *Signer) Sign(userID, email, name, provider string, ttl time.Duration) (string, error) {
	now := time.Now()
	c := claims{
		Sub:      userID,
		Email:    email,
		Name:     name,
		Provider: provider,
		Exp:      now.Add(ttl).Unix(),
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshaling header: %w", err)
	}

	claimsJSON, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("marshaling claims: %w", err)
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	message := headerB64 + "." + claimsB64

	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return message + "." + signature, nil
}

func (v *Verifier) Verify(token string) (*auth.TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format: %w", domain.ErrUnauthorized)
	}

	headerB64, claimsB64, signatureB64 := parts[0], parts[1], parts[2]

	message := headerB64 + "." + claimsB64

	mac := hmac.New(sha256.New, []byte(v.secret))
	mac.Write([]byte(message))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signatureB64), []byte(expectedSig)) {
		return nil, fmt.Errorf("invalid signature: %w", domain.ErrUnauthorized)
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(claimsB64)
	if err != nil {
		return nil, fmt.Errorf("decoding claims: %w", domain.ErrUnauthorized)
	}

	var c claims
	if err := json.Unmarshal(claimsJSON, &c); err != nil {
		return nil, fmt.Errorf("unmarshaling claims: %w", domain.ErrUnauthorized)
	}

	if time.Now().Unix() > c.Exp {
		return nil, fmt.Errorf("token expired: %w", domain.ErrUnauthorized)
	}

	return &auth.TokenClaims{
		UserID:   c.Sub,
		Email:    c.Email,
		Name:     c.Name,
		Provider: c.Provider,
	}, nil
}
