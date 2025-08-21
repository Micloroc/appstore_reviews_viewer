package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"appstorereviewsviewer/internal/application/addapp"
	"appstorereviewsviewer/internal/application/getrecentreviews"
	"appstorereviewsviewer/internal/application/reloadreviews"
	"appstorereviewsviewer/internal/domain/app"
	"appstorereviewsviewer/internal/domain/review"
	"appstorereviewsviewer/internal/infrastructure/cron"
	infrahttp "appstorereviewsviewer/internal/infrastructure/http"
	persistenceapp "appstorereviewsviewer/internal/infrastructure/persistence/app"
	persistencereview "appstorereviewsviewer/internal/infrastructure/persistence/review"
)

func main() {
	dataDir := "data"
	port := "8080"

	repos, err := setupRepositories(dataDir)
	if err != nil {
		log.Fatalf("Failed to setup repositories: %v", err)
	}

	useCases := setupUseCases(repos)
	server := infrahttp.NewServer(useCases.getRecentReviews, useCases.addApp, port)
	server.Start()

	reloadReviews := cron.NewReloadReviews(useCases.reloadReviews)
	reloadReviews.Start()

	handleGracefulShutdown(reloadReviews)
}

type repositories struct {
	reviewFile review.Repository
	reviewRSS  review.Repository
	appFile    app.Repository
}

func setupRepositories(dataDir string) (*repositories, error) {
	reviewFileRepo, err := persistencereview.NewFileRepository(dataDir)
	if err != nil {
		return nil, err
	}

	rssReviewRepo := persistencereview.NewRSSRepository()

	appFileRepo, err := persistenceapp.NewFileRepository(dataDir)
	if err != nil {
		return nil, err
	}

	return &repositories{
		reviewFile: reviewFileRepo,
		reviewRSS:  rssReviewRepo,
		appFile:    appFileRepo,
	}, nil
}

type useCases struct {
	reloadReviews    reloadreviews.UseCase
	getRecentReviews getrecentreviews.UseCase
	addApp           addapp.UseCase
}

func setupUseCases(repos *repositories) *useCases {
	reloadReviewsUseCase := reloadreviews.NewUseCase(repos.reviewFile, repos.reviewRSS, repos.appFile)
	getRecentReviewsUseCase := getrecentreviews.NewUseCase(repos.reviewFile)
	addAppUseCase := addapp.NewUseCase(repos.appFile, reloadReviewsUseCase)

	return &useCases{
		reloadReviews:    reloadReviewsUseCase,
		getRecentReviews: getRecentReviewsUseCase,
		addApp:           addAppUseCase,
	}
}

func handleGracefulShutdown(reloadReviews *cron.ReloadReviews) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down server...")

	reloadReviews.Stop()
	log.Println("Server stopped")
}
