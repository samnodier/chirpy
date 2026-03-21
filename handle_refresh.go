package main

import (
	"net/http"
	"time"

	"github.com/samnodier/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", err)
		return
	}
	type response struct {
		Token string `json:"token"`
	}

	dbToken, err := cfg.db.GetToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found or expired", err)
		return
	}
	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create a new access token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
