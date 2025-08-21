package http

import (
	"appstorereviewsviewer/internal/application/addapp"
	"appstorereviewsviewer/internal/application/getrecentreviews"
)

type Handlers struct {
	getRecentReviewsUseCase getrecentreviews.UseCase
	addAppUseCase           addapp.UseCase
}

func NewHandlers(getRecentReviewsUseCase getrecentreviews.UseCase, addAppUseCase addapp.UseCase) *Handlers {
	return &Handlers{
		getRecentReviewsUseCase: getRecentReviewsUseCase,
		addAppUseCase:           addAppUseCase,
	}
}
