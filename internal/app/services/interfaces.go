package services

type SignatureService interface {
	CreateSignature(userName string) error
	VerifySignature() error
}