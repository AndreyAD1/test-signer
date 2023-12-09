package services

type TestAnswer struct {
	Question string
	Answer   string
}

type ExternalSignature struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}
