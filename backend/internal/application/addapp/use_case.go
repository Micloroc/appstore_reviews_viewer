package addapp

import (
	"log/slog"

	"appstorereviewsviewer/internal/application/reloadreviews"
	"appstorereviewsviewer/internal/domain/app"
)

type UseCase interface {
	Execute(appID string) error
}

type useCase struct {
	appRepo              app.Repository
	reloadReviewsUseCase reloadreviews.UseCase
}

func NewUseCase(appRepo app.Repository, reloadReviewsUseCase reloadreviews.UseCase) *useCase {
	return &useCase{appRepo: appRepo, reloadReviewsUseCase: reloadReviewsUseCase}
}

func (u *useCase) Execute(appID string) error {
	app, err := app.NewApp(appID)
	if err != nil {
		return err
	}

	err = u.appRepo.Save(app)
	if err != nil {
		return err
	}

	if err := u.reloadReviewsUseCase.Execute(); err != nil {
		slog.Error("failed to execute reload reviews", "error", err)
	}

	return nil
}
