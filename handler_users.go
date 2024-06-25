package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JoshuaTapp/gopherRSS/internal/database"
	"github.com/google/uuid"
)

func (s *APIServer) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := &parameters{}
	err := DecodeJSONBody(r, params)
	if err != nil {
		s.Logger.Debug("failed to decode JSON", "params", params)
		RespondWithError(w, 400, fmt.Sprintf("error parsing JSON body: %v", err))
		return
	}

	user, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	RespondWithJSON(w, http.StatusAccepted, user)
}

func (s *APIServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// middleware checked if user exists and is AuthN
	apiKey, ok := r.Context().Value(apiKeyContextKey).(string)
	if !ok || apiKey == "" {
		RespondWithError(w, http.StatusInternalServerError, "Missing or Invalid API key")
		return
	}

	user, err := s.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}
