package services

import "time"

type TestAnswer struct {
	Question string
	Answer   string
}

type ExternalSignature struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type StoredSignature struct {
	Answers []string `json:"answers"`
	Timestamp time.Time `json:"timestamp"`
}
