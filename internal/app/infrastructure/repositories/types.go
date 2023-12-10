package repositories

import "time"

type Signature struct {
	ID        string
	RequestID string
	UserID    string
	CreatedAt time.Time
	Answers   []TestDetails
}

type TestDetails struct {
	ID       string
	Question string
	Answer   string
}
