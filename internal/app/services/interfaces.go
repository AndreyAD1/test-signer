package services

type SignatureService interface {
	CreateSignature() error
	VerifySignature() error
}