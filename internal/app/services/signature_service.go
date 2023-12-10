package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"

	"encoding/json"
	"encoding/pem"

	"log"
	"os"

	r "github.com/AndreyAD1/test-signer/internal/app/infrastructure/repositories"
	"github.com/google/uuid"
)

type SignatureSvc struct {
	signatureRepo r.SignatureRepository
	privatekey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	key           []byte
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
	return &SignatureSvc{repo, rsaPrivateKey, rsaPublicKey, []byte("my very-very-very secret")}, nil
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

	block, err := aes.NewCipher(s.key)
	if err != nil {
		log.Printf("Error creating AES cipher: %v", err)
		return []byte{}, fmt.Errorf("Error creating AES cipher: %w", err)
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Error generating nonce: %v", err)
		return []byte{}, fmt.Errorf("Error generating nonce: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Error creating GCM: %v", err)
		return []byte{}, fmt.Errorf("Error creating GCM: %q", err)
	}
	ciphertext := aesgcm.Seal(nil, nonce, sign, nil)
	ciphertext = append(nonce, ciphertext...)
	return ciphertext, nil

	// answers := []r.TestDetails{}
	// for _, item := range testAnswers {
	// 	answers = append(
	// 		answers,
	// 		r.TestDetails{Question: item.Question, Answer: item.Answer},
	// 	)
	// }
	// dbSignature := r.Signature{
	// 	ID:        uuid.New().String(),
	// 	RequestID: requestID,
	// 	UserID:    userID,
	// 	Answers:   answers,
	// }
	// _, err = s.signatureRepo.Add(ctx, dbSignature)
	// if err != nil {
	// 	return []byte{}, fmt.Errorf("can not create a signature for %s: %w", userID, err)
	// }
	// return base64.StdEncoding.EncodeToString(signature), nil
}

func (s *SignatureSvc) VerifySignature() error { return nil }
