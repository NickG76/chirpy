package main

import (
	"net/http"
	"context"
	"database/sql"
	"encoding/json"
	
	"github.com/nickg76/chirpy/internal/database"
)


func (cfg *apiConfig) handleCreateUsr(w http.ResponseWriter, r *http.request) {
	type requestBody struct {
		Email string `json:"email"`
	}	

	var body requestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil  || body.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return 
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error with the server, ah sugar!")
		return
	}

}
