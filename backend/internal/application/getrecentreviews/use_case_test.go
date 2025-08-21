package getrecentreviews_test

import (
	"testing"
	"time"

	"appstorereviewsviewer/internal/application/getrecentreviews"
	"appstorereviewsviewer/internal/domain/review"
	reviewmocks "appstorereviewsviewer/mocks/domain/review"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetRecentReviewsUseCaseTestSuite struct {
	suite.Suite
	mockReviewRepo *reviewmocks.Repository
	useCase        getrecentreviews.UseCase
}

func (s *GetRecentReviewsUseCaseTestSuite) SetupSubTest() {
	s.mockReviewRepo = reviewmocks.NewRepository(s.T())
	s.useCase = getrecentreviews.NewUseCase(s.mockReviewRepo)
}

func (s *GetRecentReviewsUseCaseTestSuite) TestExecute() {
	s.Run("should return reviews when found", func() {
		appID := "12345"
		expectedReviews := []*review.Review{
			{
				ID:          "review1",
				AppID:       appID,
				Author:      "John Doe",
				Content:     "Great app!",
				Score:       5,
				SubmittedAt: time.Now().Add(-12 * time.Hour),
				RetrievedAt: time.Now(),
			},
			{
				ID:          "review2",
				AppID:       appID,
				Author:      "Jane Smith",
				Content:     "Good app",
				Score:       4,
				SubmittedAt: time.Now().Add(-6 * time.Hour),
				RetrievedAt: time.Now(),
			},
		}

		s.mockReviewRepo.EXPECT().FindByAppIDSince(appID, mock.AnythingOfType("time.Time")).Return(expectedReviews, nil)

		reviews, err := s.useCase.Execute(appID)

		s.NoError(err)
		s.Equal(expectedReviews, reviews)
	})

	s.Run("should return empty slice when no reviews found", func() {
		appID := "12345"
		expectedReviews := []*review.Review{}

		s.mockReviewRepo.EXPECT().FindByAppIDSince(appID, mock.AnythingOfType("time.Time")).Return(expectedReviews, nil)

		reviews, err := s.useCase.Execute(appID)

		s.NoError(err)
		s.Equal(expectedReviews, reviews)
	})

	s.Run("should return error when repository fails", func() {
		appID := "12345"

		s.mockReviewRepo.EXPECT().FindByAppIDSince(appID, mock.AnythingOfType("time.Time")).Return(nil, assert.AnError)

		reviews, err := s.useCase.Execute(appID)

		s.Error(err)
		s.Nil(reviews)
	})
}

func TestGetRecentReviewsUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetRecentReviewsUseCaseTestSuite))
}
