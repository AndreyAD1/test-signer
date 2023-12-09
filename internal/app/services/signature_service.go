package services

type SignatureSvc struct {}

func NewSignatureSvc() *SignatureSvc {
	return &SignatureSvc{}
}

func (s *SignatureSvc) CreateSignature(userName string) error {return nil}

func (s *SignatureSvc) VerifySignature() error {return nil}
