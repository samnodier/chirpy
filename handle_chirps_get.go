package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samnodier/chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpsGet(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error

	if authorIDString != "" {
		authorID, err := uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Can't parse the chirp format", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsForUser(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.db.ListChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, databaseChirpToChirp(dbChirp))
	}
	respondWithJSON(w, 200, chirps)
}
