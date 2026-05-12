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
}
