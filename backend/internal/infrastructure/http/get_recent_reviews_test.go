package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"appstorereviewsviewer/internal/domain/review"
	infrahttp "appstorereviewsviewer/internal/infrastructure/http"
	addappmocks "appstorereviewsviewer/mocks/application/addapp"
	getrecentreviewsmocks "appstorereviewsviewer/mocks/application/getrecentreviews"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetRecentReviewsHandlerTestSuite struct {
	suite.Suite
	mockAddAppUseCase           *addappmocks.UseCase
	mockGetRecentReviewsUseCase *getrecentreviewsmocks.UseCase
	handlers                    *infrahttp.Handlers
}

func (s *GetRecentReviewsHandlerTestSuite) SetupSubTest() {
	s.mockAddAppUseCase = addappmocks.NewUseCase(s.T())
	s.mockGetRecentReviewsUseCase = getrecentreviewsmocks.NewUseCase(s.T())
	s.handlers = infrahttp.NewHandlers(s.mockGetRecentReviewsUseCase, s.mockAddAppUseCase)
}

func (s *GetRecentReviewsHandlerTestSuite) TestGetRecentReviews() {
	s.Run("should return reviews when valid app ID provided", func() {
		appID := "12345"
		now := time.Now()
		expectedReviews := []*review.Review{
			{
				ID:          "review1",
				AppID:       appID,
				Author:      "John Doe",
				Content:     "Great app!",
				Score:       5,
				SubmittedAt: now.Add(-2 * time.Hour),
				RetrievedAt: now,
			},
			{
				ID:          "review2",
				AppID:       appID,
				Author:      "Jane Smith",
				Content:     "Good app",
				Score:       4,
				SubmittedAt: now.Add(-1 * time.Hour),
				RetrievedAt: now,
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/app/"+appID+"/reviews/recent", nil)
		rr := httptest.NewRecorder()

		s.mockGetRecentReviewsUseCase.EXPECT().Execute(appID).Return(expectedReviews, nil)

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
		s.Equal("Content-Type", rr.Header().Get("Access-Control-Allow-Headers"))

		var response infrahttp.ReviewsResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		s.NoError(err)
		s.Len(response.Reviews, 2)

		s.Equal("review1", response.Reviews[0].ID)
		s.Equal(appID, response.Reviews[0].AppID)
		s.Equal("John Doe", response.Reviews[0].Author)
		s.Equal("Great app!", response.Reviews[0].Content)
		s.Equal(5, response.Reviews[0].Score)

		s.Equal("review2", response.Reviews[1].ID)
		s.Equal("Jane Smith", response.Reviews[1].Author)
		s.Equal("Good app", response.Reviews[1].Content)
		s.Equal(4, response.Reviews[1].Score)
	})

	s.Run("should return empty array when no reviews found", func() {
		appID := "12345"
		req := httptest.NewRequest(http.MethodGet, "/api/v1/app/"+appID+"/reviews/recent", nil)
		rr := httptest.NewRecorder()

		s.mockGetRecentReviewsUseCase.EXPECT().Execute(appID).Return([]*review.Review{}, nil)

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		s.Equal("application/json", rr.Header().Get("Content-Type"))

		var response infrahttp.ReviewsResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		s.NoError(err)
		s.Len(response.Reviews, 0)
	})

	s.Run("should return empty array when use case returns nil", func() {
		appID := "12345"
		req := httptest.NewRequest(http.MethodGet, "/api/v1/app/"+appID+"/reviews/recent", nil)
		rr := httptest.NewRecorder()

		s.mockGetRecentReviewsUseCase.EXPECT().Execute(appID).Return(nil, nil)

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		s.Equal("application/json", rr.Header().Get("Content-Type"))

		var response infrahttp.ReviewsResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		s.NoError(err)
		s.Len(response.Reviews, 0)
	})

	s.Run("should return bad request when invalid app ID in URL", func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/app//reviews/recent", nil)
		rr := httptest.NewRecorder()

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
		s.Contains(rr.Body.String(), "Invalid app ID")
	})

	s.Run("should return bad request when URL pattern does not match", func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/invalid/path", nil)
		rr := httptest.NewRecorder()

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
		s.Contains(rr.Body.String(), "Invalid app ID")
	})

	s.Run("should return internal server error when use case fails", func() {
		appID := "12345"
		req := httptest.NewRequest(http.MethodGet, "/api/v1/app/"+appID+"/reviews/recent", nil)
		rr := httptest.NewRecorder()

		s.mockGetRecentReviewsUseCase.EXPECT().Execute(appID).Return(nil, assert.AnError)

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusInternalServerError, rr.Code)
		s.Contains(rr.Body.String(), "assert.AnError general error for testing")
	})

	s.Run("should format submitted time correctly", func() {
		appID := "12345"
		submittedAt := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
		expectedReviews := []*review.Review{
			{
				ID:          "review1",
				AppID:       appID,
				Author:      "John Doe",
				Content:     "Great app!",
				Score:       5,
				SubmittedAt: submittedAt,
				RetrievedAt: time.Now(),
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/app/"+appID+"/reviews/recent", nil)
		rr := httptest.NewRecorder()

		s.mockGetRecentReviewsUseCase.EXPECT().Execute(appID).Return(expectedReviews, nil)

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusOK, rr.Code)

		var response infrahttp.ReviewsResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		s.NoError(err)
		s.Len(response.Reviews, 1)
		s.Equal("2025-01-01T12:00:00Z", response.Reviews[0].SubmittedAt)
	})

	s.Run("should handle special characters in app ID", func() {
		appID := "app-123_test"
		req := httptest.NewRequest(http.MethodGet, "/api/v1/app/"+appID+"/reviews/recent", nil)
		rr := httptest.NewRecorder()

		s.mockGetRecentReviewsUseCase.EXPECT().Execute(appID).Return([]*review.Review{}, nil)

		s.handlers.GetRecentReviews(rr, req)

		s.Equal(http.StatusOK, rr.Code)
	})
}

func TestGetRecentReviewsHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetRecentReviewsHandlerTestSuite))
}
