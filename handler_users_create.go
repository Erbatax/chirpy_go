package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/erbatax/chirpy_go/internal/auth"
	"github.com/erbatax/chirpy_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("Failed to hash password", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashedPassword})
	if err != nil {
		log.Println("Failed to create user", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	resp := response{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}
