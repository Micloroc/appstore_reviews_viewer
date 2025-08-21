package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"appstorereviewsviewer/internal/domain/app"
)

type FileRepository struct {
	dataDir string
}

type AppData struct {
	ID string `json:"id"`
}

func NewFileRepository(dataDir string) (*FileRepository, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &FileRepository{
		dataDir: dataDir,
	}, nil
}

func (r *FileRepository) FindAll() ([]*app.App, error) {
	filePath := r.getFilePath()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*app.App{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var appsData []AppData
	if err := json.Unmarshal(data, &appsData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal apps: %w", err)
	}

	var apps []*app.App
	for _, appData := range appsData {
		app, err := app.NewApp(appData.ID)
		if err != nil {
			continue
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *FileRepository) Save(app *app.App) error {
	if app == nil {
		return fmt.Errorf("app cannot be nil")
	}

	filePath := r.getFilePath()

	var existingApps []AppData
	if data, err := os.ReadFile(filePath); err == nil {
		if err := json.Unmarshal(data, &existingApps); err != nil {
			return fmt.Errorf("failed to unmarshal existing apps: %w", err)
		}
	}

	appMap := make(map[string]AppData)
	for _, existingApp := range existingApps {
		appMap[existingApp.ID] = existingApp
	}

	appData := AppData{
		ID: app.ID,
	}
	appMap[app.ID] = appData

	var allApps []AppData
	for _, app := range appMap {
		allApps = append(allApps, app)
	}

	data, err := json.MarshalIndent(allApps, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal apps: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (r *FileRepository) getFilePath() string {
	return filepath.Join(r.dataDir, "apps.json")
}
