package services

import "context"

type SignatureService interface {
	CreateSignature(context.Context, string, string, []TestAnswer) ([]byte, error)
	VerifySignature(context.Context, string, []byte) error
}
