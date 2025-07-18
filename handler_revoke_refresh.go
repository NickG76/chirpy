package main

import (
	"net/http"
	"time"
	"strings"
	"errors"
	"github.com/nickg76/chirpy/internal/database"
	"database/sql"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	refreshToken, err := parseBearerToken(authHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing Authorization header", err)
		return
	}

	// Query the refresh token from DB ---------------------  
	dbToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Check expiration or revoked
	if time.Now().After(dbToken.ExpiresAt) || dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired or revoked", err)
		return 
	}

	// Issue new JWT 
	token, err := cfg.auth.MakeJWT(dbToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})


}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	refreshToken, err := parseBearerToken(authHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing Authorization header", err)
		return
	}

	now := time.Now()
	err = cfg.db.RevokeRefreshToken(r. Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: now,
		Token:     refreshToken,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseBearerToken(header string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", errors.New("missing Bearer prefix")
	}
	return strings.TrimPrefix(header, prefix), nil
}
//---------------------------------- TOTO -------------------------------- 
