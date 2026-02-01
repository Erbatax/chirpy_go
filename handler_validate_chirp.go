package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
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
