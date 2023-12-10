package services

import "errors"

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrDuplicatedSignature = errors.New("signature already exists")
)
