package app

import "errors"

type App struct {
	ID string
}

func NewApp(id string) (*App, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return &App{ID: id}, nil
}
