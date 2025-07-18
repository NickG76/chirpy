package main

import (
	"net/http"
	"github.com/nickg76/chirpy/internal/auth"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) { 
	type loginParams struct {
		Email	 			string `json:"email"`
		Password 			string `json:"password"`
		ExpiresInSeconds	*int   `json:"expires_in_seconds"`
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
	expiresIn := time.Hour
	if params.ExpiresInSeconds != nil {
		requested := time.Duration(*params.ExpiresInSeconds) * time.Second
		if requested <= time.Hour {
			expiresIn = requested
		}
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}
	


	type loginResponse struct {
		ID 			 uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email 		 string    `json:"email"`
		Token 		 string    `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID: 	   user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
} 
