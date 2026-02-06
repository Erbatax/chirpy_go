package main

import (
	"log"
	"net/http"
	"time"

	"github.com/erbatax/chirpy_go/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("Failed to get refresh token", http.StatusUnauthorized, err)

		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	dbRefreshToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		log.Println("Failed to get refresh token from database", http.StatusUnauthorized, err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if dbRefreshToken.ExpiresAt.Before(time.Now()) {
		log.Println("Refresh token has expired", http.StatusUnauthorized)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if dbRefreshToken.RevokedAt.Valid && dbRefreshToken.RevokedAt.Time.Before(time.Now()) {
		log.Println("Refresh token has been revoked", http.StatusUnauthorized)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	token, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.jwtSecret, time.Duration(cfg.accessTokenExpiresInSeconds)*time.Second)
	if err != nil {
		log.Println("Failed to create JWT", http.StatusInternalServerError, err)

		respondWithError(w, http.StatusInternalServerError, "Failed to create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
