package addapp_test

import (
	"testing"

	"appstorereviewsviewer/internal/application/addapp"
	"appstorereviewsviewer/internal/domain/app"
	reloadreviewsmocks "appstorereviewsviewer/mocks/application/reloadreviews"
	appmocks "appstorereviewsviewer/mocks/domain/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddAppUseCaseTestSuite struct {
	suite.Suite
	mockAppRepo              *appmocks.Repository
	mockReloadReviewsUseCase *reloadreviewsmocks.UseCase
	useCase                  addapp.UseCase
}

func (s *AddAppUseCaseTestSuite) SetupSubTest() {
	s.mockAppRepo = appmocks.NewRepository(s.T())
	s.mockReloadReviewsUseCase = reloadreviewsmocks.NewUseCase(s.T())
	s.useCase = addapp.NewUseCase(s.mockAppRepo, s.mockReloadReviewsUseCase)
}

func (s *AddAppUseCaseTestSuite) TestExecute() {
	s.Run("should save app and reload reviews when valid app ID provided", func() {
		expectedApp, _ := app.NewApp("12345")
		s.mockAppRepo.EXPECT().Save(expectedApp).Return(nil)
		s.mockReloadReviewsUseCase.EXPECT().Execute().Return(nil)

		err := s.useCase.Execute("12345")

		s.NoError(err)
	})

	s.Run("should return error when app creation fails", func() {
		err := s.useCase.Execute("")

		s.Error(err)
		s.Equal("id is required", err.Error())
	})

	s.Run("should return error when app repository save fails", func() {
		expectedApp, _ := app.NewApp("12345")
		s.mockAppRepo.EXPECT().Save(expectedApp).Return(assert.AnError)

		err := s.useCase.Execute("12345")

		s.Error(err)
	})

	s.Run("should continue when reload reviews fails", func() {
		expectedApp, _ := app.NewApp("12345")
		s.mockAppRepo.EXPECT().Save(expectedApp).Return(nil)
		s.mockReloadReviewsUseCase.EXPECT().Execute().Return(assert.AnError)

		err := s.useCase.Execute("12345")

		s.NoError(err)
	})
}

func TestAddAppUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AddAppUseCaseTestSuite))
}
