package app

type Repository interface {
	FindAll() ([]*App, error)
	Save(app *App) error
}
