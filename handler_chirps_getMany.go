package main

import (
	"net/http"
	"time"

	"github.com/erbatax/chirpy_go/internal/database"
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

	queryAuthorId := r.URL.Query().Get("author_id")
	var authorId uuid.UUID
	var err error
	if queryAuthorId != "" {
		authorId, err = uuid.Parse(queryAuthorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}
	}

	var chirps []database.Chirp
	if authorId != uuid.Nil {
		chirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorId)
	} else {
		chirps, err = cfg.db.GetManyChirps(r.Context())
	}
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
