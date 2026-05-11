package config

import (
	"io"
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
	FrontendURL         string
}

func LoadConfig() *Config {
	runningEnvironment := os.Getenv("APP_ENV")
	if runningEnvironment == "" {
		runningEnvironment = "local"
	}
	if runningEnvironment == "local" {
		err := loadENVFromFile(".env.local")
		if err != nil {
			panic("failed to load .env.local file: " + err.Error())
		}
	}

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
		googleRedirectURL = "http://localhost:8081/auth/callback/google"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
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
		FrontendURL:         frontendURL,
	}
}

func loadENVFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var content []byte
	content, err = io.ReadAll(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	return nil
}
