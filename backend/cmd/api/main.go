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
	"wc2026/internal/adapters/ws"
	"wc2026/internal/application/auth"
	leaderboardApp "wc2026/internal/application/leaderboard"
	"wc2026/internal/application/match"
	predictionApp "wc2026/internal/application/prediction"
	"wc2026/internal/config"
	"wc2026/internal/infrastructure/football"
	"wc2026/internal/infrastructure/google"
	"wc2026/internal/infrastructure/jwt"
	"wc2026/internal/infrastructure/postgres"
	"wc2026/internal/middleware"
	"wc2026/migrations"

	httpAdapters "wc2026/internal/adapters/http"
)

func main() {
	cfg := config.LoadConfig()

	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := postgres.Up(db, migrations.FS); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := postgres.NewUserRepo(db)
	matchRepo := postgres.NewMatchRepo(db)
	predictionRepo := postgres.NewPredictionRepo(db)
	leaderboardRepo := postgres.NewLeaderboardRepo(db)
	tournamentRepo := postgres.NewTournamentRepo(db)

	googleProvider := google.NewProvider(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
		cfg.AllowedEmailDomains,
	)

	footballClient := football.NewClient(cfg.FootballAPIKey, tournamentRepo)

	tokenSigner := jwt.NewSigner(cfg.JWTSecret)
	tokenVerifier := jwt.NewVerifier(cfg.JWTSecret)

	hub := ws.NewHub()
	go hub.Run()

	authService := auth.NewService(userRepo, googleProvider, tokenSigner, cfg.JWTRefreshSecret)
	matchService := match.NewService(matchRepo, footballClient, leaderboardRepo, hub)
	predictionService := predictionApp.NewService(predictionRepo, matchRepo, leaderboardRepo)
	leaderboardService := leaderboardApp.NewService(leaderboardRepo, cfg.JWTSecret)

	authHandler := httpAdapters.NewAuthHandler(authService, cfg.FrontendURL)
	matchHandler := httpAdapters.NewMatchHandler(matchService)
	predictionHandler := httpAdapters.NewPredictionHandler(predictionService)
	leaderboardHandler := httpAdapters.NewLeaderboardHandler(leaderboardService)
	wsHandler := httpAdapters.NewWebSocketHandler(hub)
	healthHandler := httpAdapters.NewHealthHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler.Handle)

	mux.HandleFunc("GET /auth/google", authHandler.InitiateGoogleLogin)
	mux.HandleFunc("GET /auth/callback/google", authHandler.HandleGoogleCallback)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)

	mux.HandleFunc("GET /matches", matchHandler.ListMatches)
	mux.HandleFunc("GET /matches/{id}", matchHandler.GetMatch)

	mux.HandleFunc("GET /ws", wsHandler.Handle)

	mux.HandleFunc("GET /leaderboards/{id}", leaderboardHandler.GetLeaderboard)

	authMiddleware := middleware.RequireAuth(tokenVerifier)
	mux.Handle("GET /auth/me", authMiddleware(http.HandlerFunc(authHandler.GetCurrentUser)))
	mux.Handle("PUT /predictions", authMiddleware(http.HandlerFunc(predictionHandler.UpsertPrediction)))
	mux.Handle("GET /leaderboards/{lbID}/matches/{matchID}/predictions", authMiddleware(http.HandlerFunc(predictionHandler.ListPredictions)))
	mux.Handle("POST /leaderboards", authMiddleware(http.HandlerFunc(leaderboardHandler.CreateLeaderboard)))
	mux.Handle("GET /leaderboards", authMiddleware(http.HandlerFunc(leaderboardHandler.ListLeaderboards)))
	mux.Handle("POST /leaderboards/{id}/invite", authMiddleware(http.HandlerFunc(leaderboardHandler.GenerateInvite)))
	mux.Handle("POST /leaderboards/join", authMiddleware(http.HandlerFunc(leaderboardHandler.JoinLeaderboard)))

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
