package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL         string
	JWTSecret           string
	JWTRefreshSecret    string
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	FootballAPIKey      string
	AllowedEmailDomains []string
	Port                string
	CORSOrigin          string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	domainsStr := os.Getenv("ALLOWED_EMAIL_DOMAINS")
	var domains []string
	if domainsStr != "" {
		domains = strings.Split(domainsStr, ",")
		for i := range domains {
			domains[i] = strings.TrimSpace(domains[i])
		}
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "*"
	}

	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		googleRedirectURL = "http://localhost:8080/auth/callback/google"
	}

	return &Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTRefreshSecret:    os.Getenv("JWT_REFRESH_SECRET"),
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:   googleRedirectURL,
		FootballAPIKey:      os.Getenv("FOOTBALL_API_KEY"),
		AllowedEmailDomains: domains,
		Port:                port,
		CORSOrigin:          corsOrigin,
	}
}
