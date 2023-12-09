package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"

	r "github.com/AndreyAD1/test-signer/internal/app/infrastructure/repositories"
	"github.com/google/uuid"
)

type SignatureSvc struct {
	signatureRepo r.SignatureRepository
	publicKey *rsa.PublicKey
}

func NewSignatureSvc(repo r.SignatureRepository, pkeyFilePath string) *SignatureSvc {
	return &SignatureSvc{}
}

func (s *SignatureSvc) CreateSignature(
	ctx context.Context,
	requestID string,
	userID string,
	testAnswers []TestAnswer,
) ([]byte, error) {
	externalSignature := ExternalSignature{uuid.New().String(), userID}
	sign, err := json.Marshal(externalSignature)
	if err != nil {
		return []byte{}, err
	}
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, s.publicKey, sign)
	if err != nil {
		log.Fatalf("Error encrypting message: %v", err)
	}
	answers := []r.Answers{}
	for _, item := range testAnswers {
		answers = append(
			answers, 
			r.Answers{Question: item.Question, Answer: item.Answer},
		)
	}
	signature := r.Signature{
		ID: uuid.New().String(),
		RequestID: requestID,
		UserID: userID,
		Answers: answers,
	}
	_, err = s.signatureRepo.Add(ctx, signature)
	if err != nil {
		return []byte{}, fmt.Errorf("can not create a signature for %s: %w", userID, err)
	}
	return ciphertext, nil
}

func (s *SignatureSvc) VerifySignature() error { return nil }
