package leaderboard

import "time"

const (
	RoleOwner  = "owner"
	RoleMember = "member"
)

type Leaderboard struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Member struct {
	ID            string    `json:"id"`
	LeaderboardID string    `json:"leaderboard_id"`
	UserID        string    `json:"user_id"`
	Role          string    `json:"role"`
	JoinedAt      time.Time `json:"joined_at"`
}

type Score struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Points int    `json:"points"`
}
