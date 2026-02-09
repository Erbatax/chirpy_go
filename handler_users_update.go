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

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("Failed to hash password", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{ID: userID, Email: params.Email, HashedPassword: hashedPassword})
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
	respondWithJSON(w, http.StatusOK, resp)
}
