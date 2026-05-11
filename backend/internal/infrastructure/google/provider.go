package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"wc2026/internal/domain"
	"wc2026/internal/domain/auth"
)

type Provider struct {
	config              *oauth2.Config
	allowedEmailDomains []string
}

func NewProvider(clientID, clientSecret, redirectURL string, allowedEmailDomains []string) *Provider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Provider{
		config:              config,
		allowedEmailDomains: allowedEmailDomains,
	}
}

func (p *Provider) Name() string {
	return "google"
}

func (p *Provider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *Provider) VerifyToken(ctx context.Context, code string) (*auth.UserInfo, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %w", err)
	}

	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("fetching user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	var googleUser struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("decoding user info: %w", err)
	}

	if err := p.checkEmailDomain(googleUser.Email); err != nil {
		return nil, err
	}

	return &auth.UserInfo{
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: googleUser.Picture,
	}, nil
}

func (p *Provider) checkEmailDomain(email string) error {
	if len(p.allowedEmailDomains) == 0 {
		return nil
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format: %w", domain.ErrUnauthorized)
	}

	emailDomain := parts[1]
	for _, allowed := range p.allowedEmailDomains {
		if emailDomain == allowed {
			return nil
		}
	}

	return fmt.Errorf("email domain %s not allowed: %w", emailDomain, domain.ErrUnauthorized)
}
