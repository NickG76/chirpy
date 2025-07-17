package main

import (
	"net/http"
	"github.com/nickg76/chirpy/internal/auth"
	"encoding/json"

	"github.com/google/uuid"
)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
		Email	 string `json:"email"`
		Password string `json:"passowrd"`
	}

	var params loginParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)	
		return
	}

	type loginResponse struct {
		Message      string    `json:"message"`
		ID 			 uuid.UUID `json:"id"`
		Email 		 string    `json:"email"`
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		Message: "Login successful",
		ID: 	 user.ID,
		Email:   user.Email,
	})
} 
