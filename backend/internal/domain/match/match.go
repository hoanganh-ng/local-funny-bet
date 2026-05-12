package match

import "time"

const (
	StatusScheduled = "scheduled"
	StatusLive      = "live"
	StatusFinished  = "finished"
)

type Match struct {
	ID           string    `json:"id"`
	TournamentID string    `json:"tournament_id"`
	HomeTeam     string    `json:"home_team"`
	AwayTeam     string    `json:"away_team"`
	HomeScore    *int      `json:"home_score"`
	AwayScore    *int      `json:"away_score"`
	KickoffAt    time.Time `json:"kickoff_at"`
	Status       string    `json:"status"`
	ExternalID   *string   `json:"external_id,omitempty"`
}
