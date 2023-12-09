package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (h HandlerContainer) SignAnswersHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		var requestInfo SignAnswersRequest
		if err := json.Unmarshal(requestBody, &requestInfo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if requestInfo.ID == "" || requestInfo.JwtToken == "" || len(requestInfo.TestAnswers) == 0 {
			http.Error(
				w,
				"an empty key: required keys: 'id', 'jwt', 'test'",
				http.StatusBadRequest,
			)
			return
		}

		token, err := jwt.ParseWithClaims(
			requestInfo.JwtToken,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(
						"Unexpected signing method: %v",
						token.Header["alg"],
					)
				}

				return []byte(h.apiSecret), nil
			},
		)
		if err != nil {
			log.Printf("wrong JWT algorithm: %s", requestInfo)
			http.Error(w, "the unexpected JWT token", http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(JWTClaims)
		if !ok {
			http.Error(w, "the unexpected JWT token", http.StatusBadRequest)
			return
		}
		if claims.UserID == "" {
			log.Printf("no user ID in jwt key: %s", requestInfo)
			http.Error(w, "the unexpected JWT token", http.StatusBadRequest)
			return
		}

		err = h.SignatureSvc.CreateSignature(claims.UserID)
		if err != nil {
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

func (h HandlerContainer) VerifySignatureHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.SignatureSvc.VerifySignature()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
