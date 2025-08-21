package review

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"appstorereviewsviewer/internal/domain/review"
)

type FileRepository struct {
	dataDir string
}

type ReviewData struct {
	ID          string    `json:"id"`
	AppID       string    `json:"app_id"`
	Author      string    `json:"author"`
	Content     string    `json:"content"`
	Score       int       `json:"score"`
	SubmittedAt time.Time `json:"submitted_at"`
	RetrievedAt time.Time `json:"retrieved_at"`
}

func NewFileRepository(dataDir string) (*FileRepository, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &FileRepository{
		dataDir: dataDir,
	}, nil
}

func (r *FileRepository) FindByAppIDSince(appID string, since time.Time) ([]*review.Review, error) {
	filePath := r.getFilePath(appID)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*review.Review{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var reviewsData []ReviewData
	if err := json.Unmarshal(data, &reviewsData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reviews: %w", err)
	}

	var filteredReviews []*review.Review
	for _, reviewData := range reviewsData {
		if reviewData.SubmittedAt.After(since) || reviewData.SubmittedAt.Equal(since) {
			review := &review.Review{
				ID:          reviewData.ID,
				AppID:       reviewData.AppID,
				Author:      reviewData.Author,
				Content:     reviewData.Content,
				Score:       reviewData.Score,
				SubmittedAt: reviewData.SubmittedAt,
				RetrievedAt: reviewData.RetrievedAt,
			}
			filteredReviews = append(filteredReviews, review)
		}
	}

	return filteredReviews, nil
}

func (r *FileRepository) Save(reviews ...*review.Review) error {
	if len(reviews) == 0 {
		return nil
	}

	appID := reviews[0].AppID
	filePath := r.getFilePath(appID)

	var existingReviews []ReviewData
	if data, err := os.ReadFile(filePath); err == nil {
		if err := json.Unmarshal(data, &existingReviews); err != nil {
			return fmt.Errorf("failed to unmarshal existing reviews: %w", err)
		}
	}

	reviewMap := make(map[string]ReviewData)
	for _, review := range existingReviews {
		reviewMap[review.ID] = review
	}

	for _, review := range reviews {
		reviewData := ReviewData{
			ID:          review.ID,
			AppID:       review.AppID,
			Author:      review.Author,
			Content:     review.Content,
			Score:       review.Score,
			SubmittedAt: review.SubmittedAt,
			RetrievedAt: review.RetrievedAt,
		}
		reviewMap[review.ID] = reviewData
	}

	var allReviews []ReviewData
	for _, review := range reviewMap {
		allReviews = append(allReviews, review)
	}

	data, err := json.MarshalIndent(allReviews, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reviews: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (r *FileRepository) getFilePath(appID string) string {
	return filepath.Join(r.dataDir, fmt.Sprintf("%s_reviews.json", appID))
}
