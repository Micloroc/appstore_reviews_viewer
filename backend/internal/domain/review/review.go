package review

import "time"

const RecentReviewHourThreshold = 24 * 2

type Review struct {
	ID          string
	AppID       string
	Author      string
	Content     string
	Score       int
	SubmittedAt time.Time
	RetrievedAt time.Time
}
