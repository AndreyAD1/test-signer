package handlers

import "github.com/AndreyAD1/test-signer/internal/app/services"

type HandlerContainer struct {
	SignatureSvc *services.SignatureSvc
}