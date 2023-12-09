package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (h HandlerContainer) SignAnswersHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		testSignature, err := h.SignatureSvc.CreateSignature(claims.UserID)
		if err != nil {
			log.Printf("create signature error for %s: %s", claims.UserID, err)
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := SignAnswersResponse{Signature: string(testSignature)}
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
		err := h.SignatureSvc.VerifySignature()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
