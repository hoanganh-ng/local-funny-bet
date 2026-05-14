package prediction

import "time"

const (
	ValueHomeWin = "home"
	ValueDraw    = "draw"
	ValueAwayWin = "away"
)

type Prediction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	MatchID   string    `json:"match_id"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
	// populated by ListByMatch only, not persisted
	UserName      string  `json:"-"`
	UserAvatarURL *string `json:"-"`
}

// PredictionWithMatch is returned by ListByUser for the history endpoint.
type PredictionWithMatch struct {
	ID        string
	UserID    string
	MatchID   string
	Value     string
	UpdatedAt time.Time
	HomeTeam  string
	AwayTeam  string
	KickoffAt time.Time
	Status    string
	HomeScore *int
	AwayScore *int
}
