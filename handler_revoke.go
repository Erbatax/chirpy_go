package main

import (
	"net/http"

	"github.com/erbatax/chirpy_go/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cfg.db.RevokeRefreshToken(r.Context(), refreshToken)

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
