package review

import "time"

type Repository interface {
	FindByAppIDSince(appID string, since time.Time) ([]*Review, error)
	Save(reviews ...*Review) error
}
