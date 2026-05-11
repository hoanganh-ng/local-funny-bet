package auth

import "context"

// TODO(global): Add more OAuth providers when going global
type Provider interface {
	Name() string
	VerifyToken(ctx context.Context, token string) (*UserInfo, error)
}

type UserInfo struct {
	Email     string
	Name      string
	AvatarURL string
}
