package handlers

import (
	"github.com/AndreyAD1/test-signer/internal/app/services"
	"github.com/golang-jwt/jwt/v5"
)

type HandlerContainer struct {
	ApiSecret    string
	SignatureSvc *services.SignatureSvc
}

type SignAnswersRequest struct {
	ID          string   `json:"id"`
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
