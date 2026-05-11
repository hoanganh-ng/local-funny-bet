package scheduler

import (
	"context"
	"log"
	"time"
)

type MatchService interface {
	FetchAndStore(ctx context.Context) error
}

type Scheduler struct {
	matchService MatchService
	interval     time.Duration
	stopCh       chan struct{}
}

func NewScheduler(matchService MatchService, interval time.Duration) *Scheduler {
	return &Scheduler{
		matchService: matchService,
		interval:     interval,
		stopCh:       make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("scheduler started, polling every %v", s.interval)

	for {
		select {
		case <-ticker.C:
			if err := s.matchService.FetchAndStore(ctx); err != nil {
				log.Printf("error fetching and storing matches: %v", err)
			} else {
				log.Println("successfully fetched and stored matches")
			}
		case <-s.stopCh:
			log.Println("scheduler stopped")
			return
		case <-ctx.Done():
			log.Println("scheduler context cancelled")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}
