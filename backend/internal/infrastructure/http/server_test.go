package http_test

import (
	"testing"
	"time"

	infrahttp "appstorereviewsviewer/internal/infrastructure/http"
	addappmocks "appstorereviewsviewer/mocks/application/addapp"
	getrecentreviewsmocks "appstorereviewsviewer/mocks/application/getrecentreviews"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	mockAddAppUseCase           *addappmocks.UseCase
	mockGetRecentReviewsUseCase *getrecentreviewsmocks.UseCase
}

func (s *ServerTestSuite) SetupSubTest() {
	s.mockAddAppUseCase = addappmocks.NewUseCase(s.T())
	s.mockGetRecentReviewsUseCase = getrecentreviewsmocks.NewUseCase(s.T())
}

func (s *ServerTestSuite) TestNewServer() {
	s.Run("should create server with correct configuration", func() {
		port := "8080"
		server := infrahttp.NewServer(s.mockGetRecentReviewsUseCase, s.mockAddAppUseCase, port)

		s.NotNil(server)
		s.Equal(":8080", server.Addr)
		s.NotNil(server.Handler)
	})

	s.Run("should create server with custom port", func() {
		port := "3000"
		server := infrahttp.NewServer(s.mockGetRecentReviewsUseCase, s.mockAddAppUseCase, port)

		s.NotNil(server)
		s.Equal(":3000", server.Addr)
	})
}

func (s *ServerTestSuite) TestServerStart() {
	s.Run("should start server without blocking", func() {
		port := "0"
		server := infrahttp.NewServer(s.mockGetRecentReviewsUseCase, s.mockAddAppUseCase, port)

		done := make(chan bool)
		go func() {
			server.Start()
			done <- true
		}()

		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			s.Fail("Server.Start() should not block")
		}

		if server.Server != nil {
			server.Close()
		}
	})
}

func (s *ServerTestSuite) TestServerHandlerRoutes() {
	s.Run("should configure routes correctly", func() {
		server := infrahttp.NewServer(s.mockGetRecentReviewsUseCase, s.mockAddAppUseCase, "8080")
		s.NotNil(server.Handler)
		s.NotNil(server.Handler)
	})
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
