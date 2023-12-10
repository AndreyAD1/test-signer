package repositories

import (
	"time"

	"github.com/google/uuid"
)

type Signature struct {
	ID        uuid.UUID
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
