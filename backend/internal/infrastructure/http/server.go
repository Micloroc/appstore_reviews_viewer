package http

import (
	"log"
	"net/http"

	"appstorereviewsviewer/internal/application/addapp"
	"appstorereviewsviewer/internal/application/getrecentreviews"
)

type Server struct {
	*http.Server
}

func NewServer(getRecentReviewsUseCase getrecentreviews.UseCase, addAppUseCase addapp.UseCase, port string) *Server {
	handlers := NewHandlers(getRecentReviewsUseCase, addAppUseCase)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/app/", handlers.GetRecentReviews)
	mux.HandleFunc("/api/v1/app", handlers.AddApp)
	handler := CorsMiddleware(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return &Server{Server: server}
}

func (s *Server) Start() {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
}
