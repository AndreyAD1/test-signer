package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"encoding/json"

	"log"

	"github.com/AndreyAD1/test-signer/internal/app/infrastructure/repositories"
	r "github.com/AndreyAD1/test-signer/internal/app/infrastructure/repositories"
	"github.com/google/uuid"
)

type SignatureSvc struct {
	signatureRepo r.SignatureRepository
	cipher        cipher.AEAD
}

func NewSignatureSvc(
	repo r.SignatureRepository,
	key string,
) (*SignatureSvc, error) {
	if len([]byte(key)) < 32 {
		return nil, fmt.Errorf("too short key: %s", key)
	}
	block, err := aes.NewCipher([]byte(key)[:32])
	if err != nil {
		log.Printf("Error creating AES cipher: %v", err)
		return nil, fmt.Errorf("Error creating AES cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Error creating GCM: %v", err)
		return nil, fmt.Errorf("Error creating GCM: %q", err)
	}
	return &SignatureSvc{repo, aesgcm}, nil
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
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Error generating nonce: %v", err)
		return []byte{}, fmt.Errorf("Error generating nonce: %w", err)
	}
	ciphertext := s.cipher.Seal(nil, nonce, sign, nil)
	ciphertext = append(nonce, ciphertext...)

	answers := []repositories.TestDetails{}
	for _, a := range testAnswers {
		d := repositories.TestDetails{Question: a.Question, Answer: a.Answer}
		answers = append(answers, d)
	}
	storageSignature := repositories.Signature{
		ID: uuid.New(),
		RequestID: requestID,
		UserID: userID,
		CreatedAt: time.Now(),
		Answers: answers,
	}
	s.signatureRepo.Add(ctx, storageSignature)
	return ciphertext, nil
}

func (s *SignatureSvc) VerifySignature(ctx context.Context, username string, ciphered []byte) error {
	nonce, ciphered := ciphered[:12], ciphered[12:]
	decyphered, err := s.cipher.Open(nil, nonce, ciphered, nil)
	if err != nil {
		log.Printf("Error decrypting data: %v", err)
		return ErrInvalidSignature
	}
	var signature ExternalSignature
	if err := json.Unmarshal(decyphered, &signature); err != nil {
		log.Printf("can not unmarshal a decyphered signature: %v", err)
		return ErrInvalidSignature
	}
	fmt.Println(signature)
	return nil
}
