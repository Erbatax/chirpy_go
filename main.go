package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	apiCfg := &apiConfig{}

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serveMux.HandleFunc("GET /api/healthz", healthzHandler)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandler)

	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Invalid request payload", http.StatusBadRequest)

		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if params.Body == "" {
		log.Println("Missing 'body' parameter", http.StatusBadRequest)

		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) > maxChirpLength {
		log.Printf("Chirp too long: %d characters (max %d)", len(params.Body), maxChirpLength)

		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := replaceBadWords(params.Body, badWords)
	resp := response{
		CleanedBody: cleanedBody,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}

func replaceBadWords(text string, badWords []string) string {
	words := strings.Split(text, " ")
	lowerBadWordsSet := make(map[string]struct{})
	for _, bw := range badWords {
		lowerBadWordsSet[strings.ToLower(bw)] = struct{}{}
	}

	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := lowerBadWordsSet[lowerWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
