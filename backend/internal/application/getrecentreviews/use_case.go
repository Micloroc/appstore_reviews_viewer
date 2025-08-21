package getrecentreviews

import (
	"time"

	"appstorereviewsviewer/internal/domain/review"
)

type UseCase interface {
	Execute(appID string) ([]*review.Review, error)
}

type useCase struct {
	reviewRepo review.Repository
}

func NewUseCase(reviewRepo review.Repository) *useCase {
	return &useCase{
		reviewRepo: reviewRepo,
	}
}

func (s *useCase) Execute(appID string) ([]*review.Review, error) {
	since := time.Now().Add(-time.Duration(review.RecentReviewHourThreshold) * time.Hour)
	return s.reviewRepo.FindByAppIDSince(appID, since)
}
