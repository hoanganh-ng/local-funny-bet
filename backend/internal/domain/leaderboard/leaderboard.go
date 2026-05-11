package leaderboard

import "time"

const (
	RoleOwner  = "owner"
	RoleMember = "member"
)

type Leaderboard struct {
	ID        string
	Name      string
	CreatedBy string
	CreatedAt time.Time
}

type Member struct {
	ID            string
	LeaderboardID string
	UserID        string
	Role          string
	JoinedAt      time.Time
}

type Score struct {
	UserID string
	Name   string
	Points int
}
