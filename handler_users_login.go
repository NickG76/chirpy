package main

import (
	"net/http"
	"github.com/nickg76/chirpy/internal/database"
	"encoding/json"
	"time"
	"database/sql"

	"github.com/google/uuid"
)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) { 
	type loginParams struct {
		Email	 			string `json:"email"`
		Password 			string `json:"password"`
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

	err = cfg.auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)	
		return
	}

	accessToken, err := cfg.auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token", err)
		return
	}
	
	refreshToken, err := cfg.auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create refresh token", err)
		return
	}

	now := time.Now()
	expiresAt := now.Add(60 * 24 * time.Hour) // 60 days

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
		RevokedAt: sql.NullTime{Valid: false},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could now save refresh token", err)
		return
	}


	type loginResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
} 
