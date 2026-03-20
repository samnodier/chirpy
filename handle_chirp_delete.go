package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samnodier/chirpy/internal/auth"
	"github.com/samnodier/chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpDelete(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't parse the chirp format", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not the author of this chirp", nil)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpUUID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Could not delete this chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
