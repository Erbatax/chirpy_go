package main

import (
	"encoding/json"
	"net/http"

	"github.com/erbatax/chirpy_go/internal/auth"
	"github.com/erbatax/chirpy_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "API key missing or invalid")
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	availableEvents := map[string]struct{}{
		"user.upgraded": {},
	}

	if _, ok := availableEvents[params.Event]; !ok {
		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}

	if params.Event == "user.upgraded" {
		userId, err := uuid.Parse(params.Data.UserId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		_, err = cfg.db.UpdateUserIsChirpyRed(r.Context(), database.UpdateUserIsChirpyRedParams{
			ID:          userId,
			IsChirpyRed: true,
		})
		if err != nil {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}
}
