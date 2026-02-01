package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Println("Failed to create user", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	resp := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}
