package reloadreviews

import (
	"log/slog"
	"time"

	"appstorereviewsviewer/internal/domain/app"
	"appstorereviewsviewer/internal/domain/review"
)

type UseCase interface {
	Execute() error
}

type useCase struct {
	localReviewRepo  review.Repository
	remoteReviewRepo review.Repository
	appRepo          app.Repository
}

func NewUseCase(localReviewRepo, remoteReviewRepo review.Repository, appRepo app.Repository) *useCase {
	return &useCase{
		localReviewRepo:  localReviewRepo,
		remoteReviewRepo: remoteReviewRepo,
		appRepo:          appRepo,
	}
}

func (s *useCase) Execute() error {
	apps, err := s.appRepo.FindAll()
	if err != nil {
		return err
	}

	for _, app := range apps {
		reviews, err := s.remoteReviewRepo.FindByAppIDSince(
			app.ID,
			time.Now().Add(-time.Duration(review.RecentReviewHourThreshold)*time.Hour),
		)
		if err != nil {
			slog.Error("error finding reviews for app", "app", app.ID, "error", err)
			continue
		}

		for _, review := range reviews {
			if err := s.localReviewRepo.Save(review); err != nil {
				slog.Error("error saving review", "review", review, "error", err)
			}
		}
	}

	return nil
}
