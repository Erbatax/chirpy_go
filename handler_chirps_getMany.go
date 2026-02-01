package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getManyChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type chripResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.db.GetManyChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respChirps := make([]chripResponse, 0, len(chirps))
	for _, chirp := range chirps {
		respChirps = append(respChirps, chripResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, respChirps)
}
