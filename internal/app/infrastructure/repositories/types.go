package repositories

import "time"

type Signature struct {
	ID        string
	RequestID string
	UserID    string
	CreatedAt time.Time
	Answers   []Answers
	Signature []byte
}

type Answers struct {
	ID       string
	Question string
	Answer   string
}
