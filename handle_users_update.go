package main

import (
	"encoding/json"
	"net/http"

	"github.com/samnodier/chirpy/internal/auth"
	"github.com/samnodier/chirpy/internal/database"
)

func (cfg *apiConfig) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid access token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	currentUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	newEmail := params.Email
	if newEmail == "" {
		newEmail = currentUser.Email
	}
	newHashedPassword := currentUser.HashedPassword
	if params.Password != "" {
		newHashedPassword, err = auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't hash the password", err)
			return
		}
	}
	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          newEmail,
		HashedPassword: newHashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update the user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, databaseUserToUser(updatedUser))
}
