package review_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"appstorereviewsviewer/internal/domain/review"
	reviewRepo "appstorereviewsviewer/internal/infrastructure/persistence/review"
	"github.com/stretchr/testify/suite"
)

type ReviewFileRepositoryTestSuite struct {
	suite.Suite
	tempDir string
	repo    *reviewRepo.FileRepository
}

func (s *ReviewFileRepositoryTestSuite) SetupSubTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "review_repo_test")
	s.Require().NoError(err)

	s.repo, err = reviewRepo.NewFileRepository(s.tempDir)
	s.Require().NoError(err)
}

func (s *ReviewFileRepositoryTestSuite) TearDownSubTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ReviewFileRepositoryTestSuite) TestNewFileRepository() {
	s.Run("should create repository with valid data directory", func() {
		tempDir, err := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)
		s.Require().NoError(err)

		repo, err := reviewRepo.NewFileRepository(tempDir)

		s.NoError(err)
		s.NotNil(repo)
		s.DirExists(tempDir)
	})

	s.Run("should create data directory if it does not exist", func() {
		nonExistentDir := filepath.Join(s.tempDir, "new_dir")

		repo, err := reviewRepo.NewFileRepository(nonExistentDir)

		s.NoError(err)
		s.NotNil(repo)
		s.DirExists(nonExistentDir)
	})
}

func (s *ReviewFileRepositoryTestSuite) TestFindByAppIDSince() {
	s.Run("should return empty slice when no reviews file exists", func() {
		since := time.Now().Add(-24 * time.Hour)

		reviews, err := s.repo.FindByAppIDSince("12345", since)

		s.NoError(err)
		s.Empty(reviews)
	})

	s.Run("should return reviews when file exists with valid data", func() {
		appID := "12345"
		now := time.Now()
		testReview := &review.Review{
			ID:          "review1",
			AppID:       appID,
			Author:      "John Doe",
			Content:     "Great app!",
			Score:       5,
			SubmittedAt: now.Add(-12 * time.Hour),
			RetrievedAt: now,
		}

		err := s.repo.Save(testReview)
		s.Require().NoError(err)

		since := now.Add(-24 * time.Hour)
		reviews, err := s.repo.FindByAppIDSince(appID, since)

		s.NoError(err)
		s.Len(reviews, 1)
		s.Equal("review1", reviews[0].ID)
		s.Equal(appID, reviews[0].AppID)
		s.Equal("John Doe", reviews[0].Author)
		s.Equal("Great app!", reviews[0].Content)
		s.Equal(5, reviews[0].Score)
	})

	s.Run("should filter reviews by time correctly", func() {
		appID := "12345"
		now := time.Now()

		oldReview := &review.Review{
			ID:          "old_review",
			AppID:       appID,
			Author:      "Old User",
			Content:     "Old review",
			Score:       3,
			SubmittedAt: now.Add(-48 * time.Hour),
			RetrievedAt: now,
		}

		newReview := &review.Review{
			ID:          "new_review",
			AppID:       appID,
			Author:      "New User",
			Content:     "New review",
			Score:       5,
			SubmittedAt: now.Add(-12 * time.Hour),
			RetrievedAt: now,
		}

		err := s.repo.Save(oldReview, newReview)
		s.Require().NoError(err)

		since := now.Add(-24 * time.Hour)
		reviews, err := s.repo.FindByAppIDSince(appID, since)

		s.NoError(err)
		s.Len(reviews, 1)
		s.Equal("new_review", reviews[0].ID)
	})

	s.Run("should include reviews submitted exactly at since time", func() {
		appID := "12345"
		now := time.Now()

		exactTimeReview := &review.Review{
			ID:          "exact_review",
			AppID:       appID,
			Author:      "Exact User",
			Content:     "Exact time review",
			Score:       4,
			SubmittedAt: now.Add(-24 * time.Hour),
			RetrievedAt: now,
		}

		err := s.repo.Save(exactTimeReview)
		s.Require().NoError(err)

		since := now.Add(-24 * time.Hour)
		reviews, err := s.repo.FindByAppIDSince(appID, since)

		s.NoError(err)
		s.Len(reviews, 1)
		s.Equal("exact_review", reviews[0].ID)
	})

	s.Run("should return empty slice when file exists but is empty", func() {
		appID := "12345"
		filePath := filepath.Join(s.tempDir, appID+"_reviews.json")
		err := os.WriteFile(filePath, []byte("[]"), 0o644)
		s.Require().NoError(err)

		since := time.Now().Add(-24 * time.Hour)
		reviews, err := s.repo.FindByAppIDSince(appID, since)

		s.NoError(err)
		s.Empty(reviews)
	})

	s.Run("should return error when file contains invalid JSON", func() {
		appID := "12345"
		filePath := filepath.Join(s.tempDir, appID+"_reviews.json")
		err := os.WriteFile(filePath, []byte("invalid json"), 0o644)
		s.Require().NoError(err)

		since := time.Now().Add(-24 * time.Hour)
		reviews, err := s.repo.FindByAppIDSince(appID, since)

		s.Error(err)
		s.Nil(reviews)
		s.Contains(err.Error(), "failed to unmarshal reviews")
	})
}

func (s *ReviewFileRepositoryTestSuite) TestSave() {
	s.Run("should save single review successfully", func() {
		appID := "12345"
		testReview := &review.Review{
			ID:          "review1",
			AppID:       appID,
			Author:      "John Doe",
			Content:     "Great app!",
			Score:       5,
			SubmittedAt: time.Now().Add(-12 * time.Hour),
			RetrievedAt: time.Now(),
		}

		err := s.repo.Save(testReview)

		s.NoError(err)
		s.FileExists(filepath.Join(s.tempDir, appID+"_reviews.json"))

		reviews, err := s.repo.FindByAppIDSince(appID, time.Now().Add(-24*time.Hour))
		s.NoError(err)
		s.Len(reviews, 1)
		s.Equal("review1", reviews[0].ID)
	})

	s.Run("should save multiple reviews successfully", func() {
		appID := "12345"
		now := time.Now()

		review1 := &review.Review{
			ID:          "review1",
			AppID:       appID,
			Author:      "John Doe",
			Content:     "Great app!",
			Score:       5,
			SubmittedAt: now.Add(-12 * time.Hour),
			RetrievedAt: now,
		}

		review2 := &review.Review{
			ID:          "review2",
			AppID:       appID,
			Author:      "Jane Smith",
			Content:     "Good app",
			Score:       4,
			SubmittedAt: now.Add(-6 * time.Hour),
			RetrievedAt: now,
		}

		err := s.repo.Save(review1, review2)

		s.NoError(err)
		reviews, err := s.repo.FindByAppIDSince(appID, time.Now().Add(-24*time.Hour))
		s.NoError(err)
		s.Len(reviews, 2)
	})

	s.Run("should do nothing when no reviews provided", func() {
		err := s.repo.Save()

		s.NoError(err)
	})

	s.Run("should update existing review when saving duplicate ID", func() {
		appID := "12345"
		now := time.Now()

		originalReview := &review.Review{
			ID:          "review1",
			AppID:       appID,
			Author:      "John Doe",
			Content:     "Great app!",
			Score:       5,
			SubmittedAt: now.Add(-12 * time.Hour),
			RetrievedAt: now,
		}

		updatedReview := &review.Review{
			ID:          "review1",
			AppID:       appID,
			Author:      "John Doe",
			Content:     "Updated review!",
			Score:       4,
			SubmittedAt: now.Add(-12 * time.Hour),
			RetrievedAt: now,
		}

		err := s.repo.Save(originalReview)
		s.NoError(err)

		err = s.repo.Save(updatedReview)
		s.NoError(err)

		reviews, err := s.repo.FindByAppIDSince(appID, time.Now().Add(-24*time.Hour))
		s.NoError(err)
		s.Len(reviews, 1)
		s.Equal("Updated review!", reviews[0].Content)
		s.Equal(4, reviews[0].Score)
	})

	s.Run("should preserve existing reviews when adding new ones", func() {
		appID := "12345"
		now := time.Now()

		existingReview := &review.Review{
			ID:          "existing",
			AppID:       appID,
			Author:      "Existing User",
			Content:     "Existing review",
			Score:       3,
			SubmittedAt: now.Add(-24 * time.Hour),
			RetrievedAt: now,
		}

		newReview := &review.Review{
			ID:          "new",
			AppID:       appID,
			Author:      "New User",
			Content:     "New review",
			Score:       5,
			SubmittedAt: now.Add(-12 * time.Hour),
			RetrievedAt: now,
		}

		err := s.repo.Save(existingReview)
		s.NoError(err)

		err = s.repo.Save(newReview)
		s.NoError(err)

		reviews, err := s.repo.FindByAppIDSince(appID, time.Now().Add(-48*time.Hour))
		s.NoError(err)
		s.Len(reviews, 2)

		reviewIDs := make(map[string]bool)
		for _, review := range reviews {
			reviewIDs[review.ID] = true
		}
		s.True(reviewIDs["existing"])
		s.True(reviewIDs["new"])
	})

	s.Run("should handle corrupted existing file gracefully", func() {
		appID := "12345"
		filePath := filepath.Join(s.tempDir, appID+"_reviews.json")
		err := os.WriteFile(filePath, []byte("invalid json"), 0o644)
		s.Require().NoError(err)

		testReview := &review.Review{
			ID:          "review1",
			AppID:       appID,
			Author:      "John Doe",
			Content:     "Great app!",
			Score:       5,
			SubmittedAt: time.Now(),
			RetrievedAt: time.Now(),
		}

		err = s.repo.Save(testReview)

		s.Error(err)
		s.Contains(err.Error(), "failed to unmarshal existing reviews")
	})
}

func TestReviewFileRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewFileRepositoryTestSuite))
}
