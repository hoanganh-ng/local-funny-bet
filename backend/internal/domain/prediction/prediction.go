package prediction

import "time"

const (
	ValueHomeWin = "home"
	ValueDraw    = "draw"
	ValueAwayWin = "away"
)

type Prediction struct {
	ID        string
	UserID    string
	MatchID   string
	Value     string
	UpdatedAt time.Time
}
