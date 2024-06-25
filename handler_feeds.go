package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/JoshuaTapp/gopherRSS/internal/database"
	"github.com/google/uuid"
)

func (s *APIServer) CreateFeedHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	params := &parameters{}
	err := DecodeJSONBody(r, params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	apiKey, ok := r.Context().Value(apiKeyContextKey).(string)
	if !ok || apiKey == "" {
		RespondWithError(w, http.StatusInternalServerError, "Missing or Invalid API key")
		return
	}

	dbUser, err := s.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user := databaseUsertoUser(dbUser)
	s.Logger.Info("creating feed", slog.String("user_id", user.ID.String()))

	dbFeed, err := s.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.URL,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
	})
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Failed to create feed: %v", err))
		return
	}
	feed := databaseFeedToFeed(dbFeed)

	// have user follow feed:
	s.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})

	RespondWithJSON(w, http.StatusAccepted, feed)
}

func (s *APIServer) GetAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := s.DB.GetAllFeeds(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	feeds := databaseFeedsToFeeds(dbFeeds)
	RespondWithJSON(w, http.StatusOK, feeds)
}
