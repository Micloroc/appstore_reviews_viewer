package review

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appstorereviewsviewer/internal/domain/review"
)

type RSSRepository struct {
	client *http.Client
}

type AppStoreResponse struct {
	Feed struct {
		Entry json.RawMessage `json:"entry"`
	} `json:"feed"`
}

type AppStoreEntry struct {
	ID struct {
		Label string `json:"label"`
	} `json:"id"`
	Author struct {
		Name struct {
			Label string `json:"label"`
		} `json:"name"`
	} `json:"author"`
	Content struct {
		Label string `json:"label"`
	} `json:"content"`
	Rating struct {
		Label string `json:"label"`
	} `json:"im:rating"`
	Updated struct {
		Label string `json:"label"`
	} `json:"updated"`
}

func NewRSSRepository() *RSSRepository {
	return &RSSRepository{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *RSSRepository) FindByAppIDSince(appID string, since time.Time) ([]*review.Review, error) {
	url := fmt.Sprintf("https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/page=1/json", appID)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RSS feed returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var appStoreResp AppStoreResponse
	if err := json.Unmarshal(body, &appStoreResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	reviews := make([]*review.Review, 0)

	if len(appStoreResp.Feed.Entry) == 0 {
		return reviews, nil
	}

	var entries []AppStoreEntry
	if err := json.Unmarshal(appStoreResp.Feed.Entry, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse entries: %w", err)
	}

	for _, entry := range entries {
		score, err := strconv.Atoi(entry.Rating.Label)
		if err != nil {
			continue
		}

		updatedTime, err := time.Parse(time.RFC3339, entry.Updated.Label)
		if err != nil {
			continue
		}

		reviewID := strings.TrimPrefix(entry.ID.Label, "https://itunes.apple.com/us/reviews/")

		reviewItem := &review.Review{
			ID:          reviewID,
			AppID:       appID,
			Author:      entry.Author.Name.Label,
			Content:     entry.Content.Label,
			Score:       score,
			SubmittedAt: updatedTime,
			RetrievedAt: time.Now(),
		}

		if reviewItem.SubmittedAt.Before(since) {
			continue
		}

		reviews = append(reviews, reviewItem)
	}

	return reviews, nil
}

func (r *RSSRepository) Save(reviews ...*review.Review) error {
	return errors.New("this repository is read-only")
}
