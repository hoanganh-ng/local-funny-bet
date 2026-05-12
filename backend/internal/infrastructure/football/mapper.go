package football

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"wc2026/internal/domain/match"
)

func mapMatches(apiMatches []apiMatch, tournamentID string) ([]*match.Match, error) {
	matches := make([]*match.Match, 0, len(apiMatches))
	for _, am := range apiMatches {
		m, err := mapMatch(am, tournamentID)
		if err != nil {
			return nil, fmt.Errorf("mapping match %d: %w", am.ID, err)
		}
		matches = append(matches, m)
	}

	return matches, nil
}

func mapMatch(am apiMatch, tournamentID string) (*match.Match, error) {
	kickoffAt, err := time.Parse(time.RFC3339, am.UTCDate)
	if err != nil {
		return nil, fmt.Errorf("parsing kickoff time %s: %w", am.UTCDate, err)
	}

	status := mapStatus(am.Status)
	externalID := strconv.Itoa(am.ID)

	return &match.Match{
		ID:           uuid.New().String(),
		TournamentID: tournamentID,
		HomeTeam:     am.HomeTeam.Name,
		AwayTeam:     am.AwayTeam.Name,
		HomeScore:    am.Score.FullTime.Home,
		AwayScore:    am.Score.FullTime.Away,
		KickoffAt:    kickoffAt,
		Status:       status,
		ExternalID:   &externalID,
	}, nil
}

func mapStatus(apiStatus string) string {
	switch apiStatus {
	case "SCHEDULED", "TIMED":
		return match.StatusScheduled
	case "IN_PLAY", "PAUSED":
		return match.StatusLive
	case "FINISHED":
		return match.StatusFinished
	default:
		return match.StatusScheduled
	}
}
