package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/samnodier/chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpsGet(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortString := r.URL.Query().Get("sort")

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
	if sortString == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}
	respondWithJSON(w, 200, chirps)
}
