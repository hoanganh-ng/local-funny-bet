package tournament

type Tournament struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Season     string  `json:"season"`
	LogoURL    *string `json:"logo_url,omitempty"`
	Status     string  `json:"status"`
	ExternalID *string `json:"external_id,omitempty"`
}
