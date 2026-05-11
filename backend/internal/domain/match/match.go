package match

import "time"

const (
	StatusScheduled = "scheduled"
	StatusLive      = "live"
	StatusFinished  = "finished"
)

type Match struct {
	ID           string
	TournamentID string
	HomeTeam     string
	AwayTeam     string
	HomeScore    *int
	AwayScore    *int
	KickoffAt    time.Time
	Status       string
	ExternalID   *string
}
