package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AndreyAD1/test-signer/internal/app/services"
	"github.com/golang-jwt/jwt/v5"
)

func (h HandlerContainer) SignAnswersHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), h.Timeout*time.Second)
		defer cancel()
		if r.Method != http.MethodPost {
			errMsg := fmt.Sprintf("Invalid method: '%s'. Expect 'POST'.", r.Method)
			http.Error(w, errMsg, http.StatusMethodNotAllowed)
			return
		}
		tokens := r.Header["Authorization"]
		if len(tokens) != 1 {
			log.Printf("unexpected tokens from %s", r.RemoteAddr)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		rawToken, ok := strings.CutPrefix(tokens[0], "Bearer ")
		if !ok {
			log.Printf("unexpected tokens from %s", r.RemoteAddr)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		token, err := jwt.ParseWithClaims(
			rawToken,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(h.ApiSecret), nil
			},
		)
		if err != nil {
			log.Printf("can not parse a JWT token: %s: %s", r.RemoteAddr, err)
			http.Error(w, "the unexpected JWT token", http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			log.Printf("unexpected claims: %s: %s", r.RemoteAddr, token.Claims)
			http.Error(w, "the unexpected JWT token", http.StatusBadRequest)
			return
		}
		if claims.UserID == "" {
			log.Printf("no user ID in jwt key: %s", r.RemoteAddr)
			http.Error(w, "the unexpected JWT token", http.StatusBadRequest)
			return
		}

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var requestInfo SignAnswersRequest
		if err := json.Unmarshal(requestBody, &requestInfo); err != nil {
			http.Error(
				w,
				fmt.Sprintf("unexpected request body: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}
		if requestInfo.ID == "" || len(requestInfo.TestAnswers) == 0 {
			http.Error(
				w,
				"an empty key: required keys: 'id', 'jwt', 'test'",
				http.StatusBadRequest,
			)
			return
		}

		testInfo := []services.TestAnswer{}
		for _, answer := range requestInfo.TestAnswers {
			internalAnswer := services.TestAnswer{
				Question: answer.Question,
				Answer:   answer.Answer,
			}
			testInfo = append(testInfo, internalAnswer)
		}

		testSignature, err := h.SignatureSvc.CreateSignature(
			ctx,
			requestInfo.ID,
			claims.UserID,
			testInfo,
		)
		if errors.Is(err, services.ErrDuplicatedSignature) {
			http.Error(
				w, 
				fmt.Sprintf("A repeated request_id '%s'", requestInfo.ID), 
				http.StatusBadRequest,
			)
			return
		}
		if err != nil {
			log.Printf("create signature error for %s: %s", claims.UserID, err)
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		base64Signature := base64.StdEncoding.EncodeToString(testSignature)
		response := SignAnswersResponse{Signature: base64Signature}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("response composition error for %s: %s", claims.UserID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h HandlerContainer) VerifySignatureHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), h.Timeout*time.Second)
		defer cancel()
		if r.Method != http.MethodPost {
			errMsg := fmt.Sprintf("Invalid method: '%s'. Expect 'POST'.", r.Method)
			http.Error(w, errMsg, http.StatusMethodNotAllowed)
			return
		}
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var requestInfo VerifyRequest
		if err := json.Unmarshal(requestBody, &requestInfo); err != nil {
			http.Error(
				w,
				fmt.Sprintf("unexpected request body: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}
		if requestInfo.UserID == "" || requestInfo.Signature == "" {
			http.Error(
				w,
				"'user_id' and 'signature' are required fields'",
				http.StatusBadRequest,
			)
			return
		}
		decodedSignature, err := base64.StdEncoding.DecodeString(requestInfo.Signature)
		if err != nil {
			errMsg := fmt.Sprintf("a signature is invalid.")
			log.Printf("a signature is invalid: %v", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		err = h.SignatureSvc.VerifySignature(ctx, requestInfo.UserID, decodedSignature)
		if errors.Is(err, services.ErrInvalidSignature) {
			http.Error(w, "Unexpected signature", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "An internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
