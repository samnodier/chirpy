package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't parse the chirp format", err)
		return
	}
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, databaseChirpToChirp(dbChirp))
}
