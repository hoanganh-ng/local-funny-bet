package prediction

import (
	"testing"
	"wc2026/internal/domain/match"
)

func TestCalculatePoint(t *testing.T) {
	homeScore := func(v int) *int { return &v }
	awayScore := func(v int) *int { return &v }

	tests := []struct {
		name            string
		match           match.Match
		predictionValue string
		want            int
	}{
		{
			name: "home win - correct prediction",
			match: match.Match{
				HomeScore: homeScore(2),
				AwayScore: awayScore(1),
			},
			predictionValue: ValueHomeWin,
			want:            1,
		},
		{
			name: "home win - wrong prediction (draw)",
			match: match.Match{
				HomeScore: homeScore(2),
				AwayScore: awayScore(1),
			},
			predictionValue: ValueDraw,
			want:            0,
		},
		{
			name: "draw - correct prediction",
			match: match.Match{
				HomeScore: homeScore(1),
				AwayScore: awayScore(1),
			},
			predictionValue: ValueDraw,
			want:            1,
		},
		{
			name: "draw - wrong prediction (away)",
			match: match.Match{
				HomeScore: homeScore(1),
				AwayScore: awayScore(1),
			},
			predictionValue: ValueAwayWin,
			want:            0,
		},
		{
			name: "away win - correct prediction",
			match: match.Match{
				HomeScore: homeScore(0),
				AwayScore: awayScore(3),
			},
			predictionValue: ValueAwayWin,
			want:            1,
		},
		{
			name: "away win - wrong prediction (home)",
			match: match.Match{
				HomeScore: homeScore(0),
				AwayScore: awayScore(3),
			},
			predictionValue: ValueHomeWin,
			want:            0,
		},
		{
			name: "match not finished - nil scores",
			match: match.Match{
				HomeScore: nil,
				AwayScore: nil,
			},
			predictionValue: ValueHomeWin,
			want:            0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePoint(tt.match, tt.predictionValue)
			if got != tt.want {
				t.Errorf("CalculatePoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
