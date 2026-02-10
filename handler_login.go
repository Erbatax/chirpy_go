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

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token,omitempty"`
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

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(cfg.accessTokenExpiresInSeconds)*time.Second)
	if err != nil {
		log.Println("Failed to create JWT", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("Failed to create refresh token", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: user.ID, ExpiresAt: time.Now().Add(time.Duration(cfg.refreshTokenExpiresInSeconds) * time.Second)})
	if err != nil {
		log.Println("Failed to save refresh token", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	resp := response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
