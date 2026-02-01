package main

import (
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset is only allowed in dev environment")
		return
	}
	cfg.fileserverHits.Store(0)

	err := cfg.db.ResetAll(r.Context())
	if err != nil {
		log.Println("Failed to reset database", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
