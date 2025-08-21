package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"appstorereviewsviewer/internal/domain/app"
	appRepo "appstorereviewsviewer/internal/infrastructure/persistence/app"

	"github.com/stretchr/testify/suite"
)

type AppFileRepositoryTestSuite struct {
	suite.Suite
	tempDir string
	repo    *appRepo.FileRepository
}

func (s *AppFileRepositoryTestSuite) SetupSubTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "app_repo_test")
	s.Require().NoError(err)

	s.repo, err = appRepo.NewFileRepository(s.tempDir)
	s.Require().NoError(err)
}

func (s *AppFileRepositoryTestSuite) TearDownSubTest() {
	os.RemoveAll(s.tempDir)
}

func (s *AppFileRepositoryTestSuite) TestNewFileRepository() {
	s.Run("should create repository with valid data directory", func() {
		tempDir, err := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)
		s.Require().NoError(err)

		repo, err := appRepo.NewFileRepository(tempDir)

		s.NoError(err)
		s.NotNil(repo)
		s.DirExists(tempDir)
	})

	s.Run("should create data directory if it does not exist", func() {
		nonExistentDir := filepath.Join(s.tempDir, "new_dir")

		repo, err := appRepo.NewFileRepository(nonExistentDir)

		s.NoError(err)
		s.NotNil(repo)
		s.DirExists(nonExistentDir)
	})
}

func (s *AppFileRepositoryTestSuite) TestFindAll() {
	s.Run("should return empty slice when no apps file exists", func() {
		apps, err := s.repo.FindAll()

		s.NoError(err)
		s.Empty(apps)
	})

	s.Run("should return apps when file exists with valid data", func() {
		testApp, _ := app.NewApp("12345")
		err := s.repo.Save(testApp)
		s.Require().NoError(err)

		apps, err := s.repo.FindAll()

		s.NoError(err)
		s.Len(apps, 1)
		s.Equal("12345", apps[0].ID)
	})

	s.Run("should return empty slice when file exists but is empty", func() {
		filePath := filepath.Join(s.tempDir, "apps.json")
		err := os.WriteFile(filePath, []byte("[]"), 0644)
		s.Require().NoError(err)

		apps, err := s.repo.FindAll()

		s.NoError(err)
		s.Empty(apps)
	})

	s.Run("should return error when file contains invalid JSON", func() {
		filePath := filepath.Join(s.tempDir, "apps.json")
		err := os.WriteFile(filePath, []byte("invalid json"), 0644)
		s.Require().NoError(err)

		apps, err := s.repo.FindAll()

		s.Error(err)
		s.Nil(apps)
		s.Contains(err.Error(), "failed to unmarshal apps")
	})

	s.Run("should skip invalid apps and continue processing", func() {
		filePath := filepath.Join(s.tempDir, "apps.json")
		invalidAppsJSON := `[{"id": "12345"}, {"id": ""}, {"id": "67890"}]`
		err := os.WriteFile(filePath, []byte(invalidAppsJSON), 0644)
		s.Require().NoError(err)

		apps, err := s.repo.FindAll()

		s.NoError(err)
		s.Len(apps, 2)
		s.Equal("12345", apps[0].ID)
		s.Equal("67890", apps[1].ID)
	})
}

func (s *AppFileRepositoryTestSuite) TestSave() {
	s.Run("should save app successfully when valid app provided", func() {
		testApp, _ := app.NewApp("12345")

		err := s.repo.Save(testApp)

		s.NoError(err)
		s.FileExists(filepath.Join(s.tempDir, "apps.json"))

		apps, err := s.repo.FindAll()
		s.NoError(err)
		s.Len(apps, 1)
		s.Equal("12345", apps[0].ID)
	})

	s.Run("should return error when app is nil", func() {
		err := s.repo.Save(nil)

		s.Error(err)
		s.Contains(err.Error(), "app cannot be nil")
	})

	s.Run("should update existing app when saving duplicate ID", func() {
		testApp1, _ := app.NewApp("12345")
		testApp2, _ := app.NewApp("12345")

		err := s.repo.Save(testApp1)
		s.NoError(err)

		err = s.repo.Save(testApp2)
		s.NoError(err)

		apps, err := s.repo.FindAll()
		s.NoError(err)
		s.Len(apps, 1)
		s.Equal("12345", apps[0].ID)
	})

	s.Run("should save multiple different apps", func() {
		testApp1, _ := app.NewApp("12345")
		testApp2, _ := app.NewApp("67890")

		err := s.repo.Save(testApp1)
		s.NoError(err)

		err = s.repo.Save(testApp2)
		s.NoError(err)

		apps, err := s.repo.FindAll()
		s.NoError(err)
		s.Len(apps, 2)

		appIDs := make(map[string]bool)
		for _, app := range apps {
			appIDs[app.ID] = true
		}
		s.True(appIDs["12345"])
		s.True(appIDs["67890"])
	})

	s.Run("should handle corrupted existing file gracefully", func() {
		filePath := filepath.Join(s.tempDir, "apps.json")
		err := os.WriteFile(filePath, []byte("invalid json"), 0644)
		s.Require().NoError(err)

		testApp, _ := app.NewApp("12345")
		err = s.repo.Save(testApp)

		s.Error(err)
		s.Contains(err.Error(), "failed to unmarshal existing apps")
	})
}

func TestAppFileRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AppFileRepositoryTestSuite))
}
