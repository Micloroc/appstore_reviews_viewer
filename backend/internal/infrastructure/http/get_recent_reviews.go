package http

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"appstorereviewsviewer/internal/domain/review"
)

type ReviewResponse struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	Score       int    `json:"score"`
	Author      string `json:"author"`
	SubmittedAt string `json:"submittedAt"`
	AppID       string `json:"appId"`
}

type ReviewsResponse struct {
	Reviews []ReviewResponse `json:"reviews"`
}

func (h *Handlers) GetRecentReviews(w http.ResponseWriter, r *http.Request) {
	appID := extractAppIDFromPath(r.URL.Path)
	if appID == "" {
		http.Error(w, "Invalid app ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.getRecentReviewsUseCase.Execute(appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reviews == nil {
		reviews = []*review.Review{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	responseReviews := make([]ReviewResponse, len(reviews))
	for i, review := range reviews {
		responseReviews[i] = ReviewResponse{
			ID:          review.ID,
			Content:     review.Content,
			Score:       review.Score,
			Author:      review.Author,
			SubmittedAt: review.SubmittedAt.Format(time.RFC3339),
			AppID:       review.AppID,
		}
	}

	response := ReviewsResponse{
		Reviews: responseReviews,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func extractAppIDFromPath(urlPath string) string {
	re := regexp.MustCompile(`^/api/v1/app/([^/]+)/reviews/recent$`)
	matches := re.FindStringSubmatch(urlPath)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
