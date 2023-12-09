package services

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	r "github.com/AndreyAD1/test-signer/internal/app/infrastructure/repositories"
	"github.com/google/uuid"
)

type SignatureSvc struct {
	signatureRepo r.SignatureRepository
	privatekey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
}

func NewSignatureSvc(
	repo r.SignatureRepository,
	privateKeyFilePath,
	pubKeyFilePath string,
) (*SignatureSvc, error) {
	raw, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(raw))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	rsaPrivateKey := parseResult.(*rsa.PrivateKey)

	raw, err = os.ReadFile(pubKeyFilePath)
	if err != nil {
		return nil, err
	}
	block, _ = pem.Decode([]byte(raw))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return &SignatureSvc{repo, rsaPrivateKey, rsaPublicKey}, nil
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

	hash := sha256.New()
	_, err = hash.Write(sign)
	if err != nil {
		log.Fatalf("Error hashing message: %v", err)
	}
	hashedMessage := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privatekey, crypto.SHA256, hashedMessage)
	if err != nil {
		log.Fatalf("Error signing message: %v", err)
	}

	answers := []r.Answers{}
	for _, item := range testAnswers {
		answers = append(
			answers,
			r.Answers{Question: item.Question, Answer: item.Answer},
		)
	}
	dbSignature := r.Signature{
		ID:        uuid.New().String(),
		RequestID: requestID,
		UserID:    userID,
		Answers:   answers,
		Signature: signature,
	}
	_, err = s.signatureRepo.Add(ctx, dbSignature)
	if err != nil {
		return []byte{}, fmt.Errorf("can not create a signature for %s: %w", userID, err)
	}
	return signature, nil
}

func (s *SignatureSvc) VerifySignature() error { return nil }
