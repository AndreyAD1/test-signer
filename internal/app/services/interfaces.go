package services

type SignatureService interface {
	CreateSignature(userName string) ([]byte, error)
	VerifySignature() error
}