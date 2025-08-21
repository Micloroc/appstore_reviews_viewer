package reloadreviews_test

import (
	"testing"
	"time"

	"appstorereviewsviewer/internal/application/reloadreviews"
	"appstorereviewsviewer/internal/domain/app"
	"appstorereviewsviewer/internal/domain/review"
	appmocks "appstorereviewsviewer/mocks/domain/app"
	reviewmocks "appstorereviewsviewer/mocks/domain/review"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ReloadReviewsUseCaseTestSuite struct {
	suite.Suite
	mockLocalReviewRepo  *reviewmocks.Repository
	mockRemoteReviewRepo *reviewmocks.Repository
	mockAppRepo          *appmocks.Repository
	useCase              reloadreviews.UseCase
}

func (s *ReloadReviewsUseCaseTestSuite) SetupSubTest() {
	s.mockLocalReviewRepo = reviewmocks.NewRepository(s.T())
	s.mockRemoteReviewRepo = reviewmocks.NewRepository(s.T())
	s.mockAppRepo = appmocks.NewRepository(s.T())
	s.useCase = reloadreviews.NewUseCase(
		s.mockLocalReviewRepo,
		s.mockRemoteReviewRepo,
		s.mockAppRepo,
	)
}

func (s *ReloadReviewsUseCaseTestSuite) TestExecute() {
	s.Run("should reload reviews for all apps successfully", func() {
		apps := []*app.App{
			{ID: "app1"},
			{ID: "app2"},
		}

		app1Reviews := []*review.Review{
			{
				ID:          "review1",
				AppID:       "app1",
				Author:      "John Doe",
				Content:     "Great app!",
				Score:       5,
				SubmittedAt: time.Now().Add(-12 * time.Hour),
				RetrievedAt: time.Now(),
			},
		}

		app2Reviews := []*review.Review{
			{
				ID:          "review2",
				AppID:       "app2",
				Author:      "Jane Smith",
				Content:     "Good app",
				Score:       4,
				SubmittedAt: time.Now().Add(-6 * time.Hour),
				RetrievedAt: time.Now(),
			},
		}

		s.mockAppRepo.EXPECT().FindAll().Return(apps, nil)
		s.mockRemoteReviewRepo.EXPECT().FindByAppIDSince("app1", mock.AnythingOfType("time.Time")).Return(app1Reviews, nil)
		s.mockRemoteReviewRepo.EXPECT().FindByAppIDSince("app2", mock.AnythingOfType("time.Time")).Return(app2Reviews, nil)
		s.mockLocalReviewRepo.EXPECT().Save(mock.Anything).Return(nil)
		s.mockLocalReviewRepo.EXPECT().Save(mock.Anything).Return(nil)

		err := s.useCase.Execute()

		s.NoError(err)
	})

	s.Run("should return error when app repository fails", func() {
		s.mockAppRepo.EXPECT().FindAll().Return(nil, assert.AnError)

		err := s.useCase.Execute()

		s.Error(err)
	})

	s.Run("should continue when remote repository fails for one app", func() {
		apps := []*app.App{
			{ID: "app1"},
			{ID: "app2"},
		}

		app2Reviews := []*review.Review{
			{
				ID:          "review2",
				AppID:       "app2",
				Author:      "Jane Smith",
				Content:     "Good app",
				Score:       4,
				SubmittedAt: time.Now().Add(-6 * time.Hour),
				RetrievedAt: time.Now(),
			},
		}

		s.mockAppRepo.EXPECT().FindAll().Return(apps, nil)
		s.mockRemoteReviewRepo.EXPECT().FindByAppIDSince("app1", mock.AnythingOfType("time.Time")).Return(nil, assert.AnError)
		s.mockRemoteReviewRepo.EXPECT().FindByAppIDSince("app2", mock.AnythingOfType("time.Time")).Return(app2Reviews, nil)
		s.mockLocalReviewRepo.EXPECT().Save(mock.Anything).Return(nil)

		err := s.useCase.Execute()

		s.NoError(err)
	})

	s.Run("should continue when local repository save fails for one review", func() {
		apps := []*app.App{
			{ID: "app1"},
		}

		app1Reviews := []*review.Review{
			{
				ID:          "review1",
				AppID:       "app1",
				Author:      "John Doe",
				Content:     "Great app!",
				Score:       5,
				SubmittedAt: time.Now().Add(-12 * time.Hour),
				RetrievedAt: time.Now(),
			},
		}

		s.mockAppRepo.EXPECT().FindAll().Return(apps, nil)
		s.mockRemoteReviewRepo.EXPECT().FindByAppIDSince("app1", mock.AnythingOfType("time.Time")).Return(app1Reviews, nil)
		s.mockLocalReviewRepo.EXPECT().Save(mock.Anything).Return(assert.AnError)

		err := s.useCase.Execute()

		s.NoError(err)
	})

	s.Run("should handle empty apps list", func() {
		apps := []*app.App{}

		s.mockAppRepo.EXPECT().FindAll().Return(apps, nil)

		err := s.useCase.Execute()

		s.NoError(err)
	})

	s.Run("should handle empty reviews for app", func() {
		apps := []*app.App{
			{ID: "app1"},
		}

		s.mockAppRepo.EXPECT().FindAll().Return(apps, nil)
		s.mockRemoteReviewRepo.EXPECT().FindByAppIDSince("app1", mock.AnythingOfType("time.Time")).Return([]*review.Review{}, nil)

		err := s.useCase.Execute()

		s.NoError(err)
	})
}

func TestReloadReviewsUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ReloadReviewsUseCaseTestSuite))
}
