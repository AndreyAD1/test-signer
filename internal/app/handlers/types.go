package handlers

import (
	"time"

	"github.com/AndreyAD1/test-signer/internal/app/services"
	"github.com/golang-jwt/jwt/v5"
)

type HandlerContainer struct {
	ApiSecret    string
	SignatureSvc services.SignatureService
	Timeout      time.Duration
}

type SignAnswersRequest struct {
	ID          string   `json:"id"` // an idempotency key
	TestAnswers []answer `json:"test"`
}

type answer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type SignAnswersResponse struct {
	Signature string `json:"signature"`
}

type VerifyRequest struct {
	UserID    string `json:"user_id"`
	Signature string `json:"signature"`
}
