package main

import (
	"net/http"
	"github.com/nickg76/chirpy/internal/auth"
)
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
		Email	 string `json:"email"`
		Password string `json:"passowrd"`
	}
} 
