package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "appstorereviewsviewer/internal/infrastructure/http"
	"github.com/stretchr/testify/suite"
)

type CorsMiddlewareTestSuite struct {
	suite.Suite
}

func (s *CorsMiddlewareTestSuite) TestCorsMiddleware() {
	s.Run("should add CORS headers to response", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("test response"))
		})

		corsHandler := infrahttp.CorsMiddleware(handler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
		s.Equal("GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		s.Equal("Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
		s.Equal("test response", rr.Body.String())
	})

	s.Run("should handle OPTIONS request and return OK", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Fail("Handler should not be called for OPTIONS request")
		})

		corsHandler := infrahttp.CorsMiddleware(handler)
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
		s.Equal("GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		s.Equal("Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
		s.Empty(rr.Body.String())
	})

	s.Run("should preserve existing headers from wrapped handler", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Custom-Header", "custom-value")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"message": "created"}`))
		})

		corsHandler := infrahttp.CorsMiddleware(handler)
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		s.Equal(http.StatusCreated, rr.Code)
		s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
		s.Equal("GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		s.Equal("Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
		s.Equal("custom-value", rr.Header().Get("Custom-Header"))
		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(`{"message": "created"}`, rr.Body.String())
	})

	s.Run("should work with different HTTP methods", func() {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

		for _, method := range methods {
			s.Run("method "+method, func() {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})

				corsHandler := infrahttp.CorsMiddleware(handler)
				req := httptest.NewRequest(method, "/test", nil)
				rr := httptest.NewRecorder()

				corsHandler.ServeHTTP(rr, req)

				s.Equal(http.StatusOK, rr.Code)
				s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
				s.Equal("GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
				s.Equal("Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
			})
		}
	})

	s.Run("should handle handler that panics gracefully", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		corsHandler := infrahttp.CorsMiddleware(handler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()

		s.Panics(func() {
			corsHandler.ServeHTTP(rr, req)
		})

		s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
		s.Equal("GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		s.Equal("Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	s.Run("should handle root path request", func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		corsHandler := infrahttp.CorsMiddleware(handler)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		corsHandler.ServeHTTP(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		s.Equal("*", rr.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestCorsMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CorsMiddlewareTestSuite))
}
