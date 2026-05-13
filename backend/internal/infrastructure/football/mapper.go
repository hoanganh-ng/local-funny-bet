package football

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"wc2026/internal/domain/match"
)

func mapMatches(apiMatches []apiMatch) []*match.Match {
	matches := make([]*match.Match, 0, len(apiMatches))
	for _, am := range apiMatches {
		m, err := mapMatch(am)
		if err != nil {
			log.Printf("skipping malformed match %d: %v", am.ID, err)
			continue
		}
		matches = append(matches, m)
	}

	return matches
}

func mapMatch(am apiMatch) (*match.Match, error) {
	kickoffAt, err := time.Parse(time.RFC3339, am.UTCDate)
	if err != nil {
		return nil, err
	}

	status := mapStatus(am.Status)
	externalID := strconv.Itoa(am.ID)

	return &match.Match{
		ID:           uuid.New().String(),
		TournamentID: "",
		HomeTeam:     am.HomeTeam.Name,
		AwayTeam:     am.AwayTeam.Name,
		HomeTeamCode: am.HomeTeam.TLA,
		AwayTeamCode: am.AwayTeam.TLA,
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
