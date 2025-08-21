package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "appstorereviewsviewer/internal/infrastructure/http"
	addappmocks "appstorereviewsviewer/mocks/application/addapp"
	getrecentreviewsmocks "appstorereviewsviewer/mocks/application/getrecentreviews"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AddAppHandlerTestSuite struct {
	suite.Suite
	mockAddAppUseCase           *addappmocks.UseCase
	mockGetRecentReviewsUseCase *getrecentreviewsmocks.UseCase
	handlers                    *infrahttp.Handlers
}

func (s *AddAppHandlerTestSuite) SetupSubTest() {
	s.mockAddAppUseCase = addappmocks.NewUseCase(s.T())
	s.mockGetRecentReviewsUseCase = getrecentreviewsmocks.NewUseCase(s.T())
	s.handlers = infrahttp.NewHandlers(s.mockGetRecentReviewsUseCase, s.mockAddAppUseCase)
}

func (s *AddAppHandlerTestSuite) TestAddApp() {
	s.Run("should add app successfully when valid request provided", func() {
		requestBody := infrahttp.AddAppRequest{AppID: "12345"}
		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/app", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		s.mockAddAppUseCase.EXPECT().Execute("12345").Return(nil)

		s.handlers.AddApp(rr, req)

		s.Equal(http.StatusCreated, rr.Code)
		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
		s.Equal("Content-Type", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	s.Run("should return method not allowed when non-POST method used", func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/app", nil)
		rr := httptest.NewRecorder()

		s.handlers.AddApp(rr, req)

		s.Equal(http.StatusMethodNotAllowed, rr.Code)
		s.Contains(rr.Body.String(), "Method not allowed")
	})

	s.Run("should return bad request when invalid JSON provided", func() {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/app", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		s.handlers.AddApp(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
		s.Contains(rr.Body.String(), "Invalid request body")
	})

	s.Run("should return bad request when empty app ID provided", func() {
		requestBody := infrahttp.AddAppRequest{AppID: ""}
		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/app", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		s.handlers.AddApp(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
		s.Contains(rr.Body.String(), "AppID is required")
	})

	s.Run("should return internal server error when use case fails", func() {
		requestBody := infrahttp.AddAppRequest{AppID: "12345"}
		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/app", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		s.mockAddAppUseCase.EXPECT().Execute("12345").Return(assert.AnError)

		s.handlers.AddApp(rr, req)

		s.Equal(http.StatusInternalServerError, rr.Code)
		s.Contains(rr.Body.String(), "assert.AnError general error for testing")
	})
}

func TestAddAppHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AddAppHandlerTestSuite))
}
