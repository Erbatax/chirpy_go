package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/erbatax/chirpy_go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int64 `json:"expires_in_seconds,omitempty"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Invalid request payload", http.StatusBadRequest)

		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if params.Email == "" {
		log.Println("Missing 'email' parameter", http.StatusBadRequest)

		respondWithError(w, http.StatusBadRequest, "Missing 'email' parameter")
		return
	}

	if params.Password == "" {
		log.Println("Missing 'password' parameter", http.StatusBadRequest)

		respondWithError(w, http.StatusBadRequest, "Missing 'password' parameter")
		return
	}

	expiresInSeconds := int64(3600) // Default to 1 hour
	if params.ExpiresInSeconds != nil {
		expiresInSeconds = int64(math.Max(0, math.Min(float64(*params.ExpiresInSeconds), float64(expiresInSeconds))))
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Println("Incorrect email", http.StatusUnauthorized, err)

		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	valid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Println("Failed to check password hash", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if !valid {
		log.Println("Incorrect password", http.StatusUnauthorized)

		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(expiresInSeconds)*time.Second)
	if err != nil {
		log.Println("Failed to create JWT", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
		return
	}

	resp := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
