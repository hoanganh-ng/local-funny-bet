package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wc2026/internal/adapters/scheduler"
	"wc2026/internal/application/auth"
	"wc2026/internal/application/match"
	"wc2026/internal/config"
	"wc2026/internal/infrastructure/football"
	"wc2026/internal/infrastructure/google"
	"wc2026/internal/infrastructure/jwt"
	"wc2026/internal/infrastructure/postgres"
	"wc2026/internal/middleware"

	httpAdapters "wc2026/internal/adapters/http"
)

func main() {
	cfg := config.LoadConfig()

	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepo(db)
	matchRepo := postgres.NewMatchRepo(db)

	googleProvider := google.NewProvider(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
		cfg.AllowedEmailDomains,
	)

	footballClient := football.NewClient(cfg.FootballAPIKey)

	tokenSigner := jwt.NewSigner(cfg.JWTSecret)

	authService := auth.NewService(userRepo, googleProvider, tokenSigner, cfg.JWTRefreshSecret)
	matchService := match.NewService(matchRepo, footballClient)

	authHandler := httpAdapters.NewAuthHandler(authService)
	matchHandler := httpAdapters.NewMatchHandler(matchService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/callback/google", authHandler.HandleGoogleCallback)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	mux.HandleFunc("GET /matches", matchHandler.ListMatches)
	mux.HandleFunc("GET /matches/{id}", matchHandler.GetMatch)

	handler := middleware.Logger(
		middleware.CORS(cfg.CORSOrigin)(mux),
	)

	matchScheduler := scheduler.NewScheduler(matchService, 60*time.Second)
	go matchScheduler.Start(context.Background())

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("shutting down gracefully...")

	matchScheduler.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
