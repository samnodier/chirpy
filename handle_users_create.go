package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/samnodier/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}

func (s *apiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	}
	user, err := s.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			respondWithError(w, 409, "User with that email already exists", pqErr)
		}
		respondWithError(w, http.StatusInternalServerError, "Error writing to database", err)
		return
	}
	respondWithJSON(w, 201, databaseUserToUser(user))
}
