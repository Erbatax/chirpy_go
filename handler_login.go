package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/erbatax/chirpy_go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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

	resp := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
