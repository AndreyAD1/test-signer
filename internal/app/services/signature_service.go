package services

type SignatureSvc struct {}

func NewSignatureSvc() *SignatureSvc {
	return &SignatureSvc{}
}

func (s *SignatureSvc) CreateSignature() error {return nil}

func (s *SignatureSvc) VerifySignature() error {return nil}
