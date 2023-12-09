package services

import (
	"context"

	r "github.com/AndreyAD1/test-signer/internal/app/infrastructure/repositories"
)

type SignatureSvc struct {
	signatureRepo r.SignatureRepository
}

func NewSignatureSvc() *SignatureSvc {
	return &SignatureSvc{}
}

func (s *SignatureSvc) CreateSignature(
	ctx context.Context,
	requestID string,
	userName string,
	testAnswers []TestAnswer,
) ([]byte, error) {

	return []byte{}, nil
}

func (s *SignatureSvc) VerifySignature() error { return nil }
