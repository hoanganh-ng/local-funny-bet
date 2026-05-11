package auth

import "context"

// TODO(global): Add more OAuth providers when going global
type Provider interface {
	Name() string
	GetAuthURL(state string) string
	VerifyToken(ctx context.Context, token string) (*UserInfo, error)
}

type UserInfo struct {
	Email     string
	Name      string
	AvatarURL string
}
