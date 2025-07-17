package main

import (
	"net/http"
	"encoding/json"
	"strings"
)



func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ahhh something went wrong!")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Ya chirp is too long")
		return 
	}

	clean := cleanChirp(req.Body)

	respondWithJson(w, http.StatusOK, responseBody{CleanedBody: clean})

}

var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func cleanChirp(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		for _, bad := range badWords {
			if strings.ToLower(word) == bad {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}
