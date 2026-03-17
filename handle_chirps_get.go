package main

import (
	"net/http"
)

func (cfg *apiConfig) handleChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirps := []Chirp{}
	c, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
	}
	for _, chirp := range c {
		chirps = append(chirps, databaseChirpToChirp(chirp))
	}
	respondWithJSON(w, 200, chirps)
}
