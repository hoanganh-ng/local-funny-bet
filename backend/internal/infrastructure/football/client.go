package football

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"wc2026/internal/domain/match"
)

const baseURL = "https://api.football-data.org/v4"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *Client) FetchMatches(ctx context.Context, competitionCode string) ([]*match.Match, error) {
	url := fmt.Sprintf("%s/competitions/%s/matches", baseURL, competitionCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching matches from football-data.org: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("football-data.org returned status %d", resp.StatusCode)
	}

	var response apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	matches, err := mapMatches(response.Matches, competitionCode)
	if err != nil {
		return nil, fmt.Errorf("mapping matches: %w", err)
	}

	return matches, nil
}

type apiResponse struct {
	Matches []apiMatch `json:"matches"`
}

type apiMatch struct {
	ID        int       `json:"id"`
	UTCDate   string    `json:"utcDate"`
	Status    string    `json:"status"`
	HomeTeam  apiTeam   `json:"homeTeam"`
	AwayTeam  apiTeam   `json:"awayTeam"`
	Score     apiScore  `json:"score"`
}

type apiTeam struct {
	Name string `json:"name"`
}

type apiScore struct {
	FullTime apiFullTime `json:"fullTime"`
}

type apiFullTime struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}
