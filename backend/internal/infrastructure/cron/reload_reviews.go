package cron

import (
	"log/slog"
	"time"

	"appstorereviewsviewer/internal/application/reloadreviews"
)

type ReloadReviews struct {
	useCase   reloadreviews.UseCase
	ticker    *time.Ticker
	stopChan  chan struct{}
	isRunning bool
}

func NewReloadReviews(useCase reloadreviews.UseCase) *ReloadReviews {
	return &ReloadReviews{
		useCase:  useCase,
		stopChan: make(chan struct{}),
	}
}

func (s *ReloadReviews) Start() {
	if s.isRunning {
		return
	}

	s.isRunning = true
	s.ticker = time.NewTicker(time.Minute)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.executeReload()
			case <-s.stopChan:
				s.ticker.Stop()
				return
			}
		}
	}()

	slog.Info("ReloadReviews started - will run every minute")
}

func (s *ReloadReviews) Stop() {
	if !s.isRunning {
		return
	}

	s.isRunning = false
	close(s.stopChan)
	slog.Info("ReloadReviews stopped")
}

func (s *ReloadReviews) executeReload() {
	if err := s.useCase.Execute(); err != nil {
		slog.Error("failed to execute reload reviews", "error", err)
	} else {
		slog.Info("reload reviews completed successfully")
	}
}
