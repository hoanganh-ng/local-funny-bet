package prediction

import "wc2026/internal/domain/match"

// CalculatePoint returns 1 if prediction matches match result, 0 otherwise.
// Returns 0 if match not finished (scores are nil).
func CalculatePoint(m match.Match, predictionValue string) int {
	if m.HomeScore == nil || m.AwayScore == nil {
		return 0
	}

	homeScore := *m.HomeScore
	awayScore := *m.AwayScore

	var actualResult string
	if homeScore > awayScore {
		actualResult = ValueHomeWin
	} else if homeScore < awayScore {
		actualResult = ValueAwayWin
	} else {
		actualResult = ValueDraw
	}

	if actualResult == predictionValue {
		return 1
	}
	return 0
}
