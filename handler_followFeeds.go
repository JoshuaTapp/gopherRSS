package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JoshuaTapp/gopherRSS/internal/database"
	"github.com/google/uuid"
)

func (s *APIServer) FollowFeedHandler(w http.ResponseWriter, r *http.Request) {
	// get user_id
	user, err := s.DB.GetUser(r.Context(), r.Context().Value(apiKeyContextKey).(string))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user_id := user.ID

	// get feed_id
	type parameters struct {
		FeedID string `json:"feed_id"`
	}
	params := &parameters{}
	err = DecodeJSONBody(r, params)
	if err != nil {
		s.Logger.Debug("failed to decode JSON", "params", params.FeedID)
		RespondWithError(w, 400, fmt.Sprintf("error parsing JSON body: %v", err))
		return
	}
	feed_id, err := uuid.Parse(params.FeedID)
	if err != nil {
		s.Logger.Debug("failed to parse UUID from feed_id", "feed_id", params.FeedID)
		RespondWithError(w, 400, "invalid feed_id")
		return
	}
	// create follow entry
	follow_rec, err := s.DB.FollowFeed(r.Context(), database.FollowFeedParams{
		FeedID:    feed_id,
		UserID:    user_id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Failed to create feed follow record: %v", err))
		return
	}

	payload := databaseFeedFollowToFeedFollow(follow_rec)
	RespondWithJSON(w, http.StatusAccepted, payload)
}

func (s *APIServer) DeleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	// get user_id
	user, err := s.DB.GetUser(r.Context(), r.Context().Value(apiKeyContextKey).(string))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user_id := user.ID

	feed_id, err := uuid.Parse(r.PathValue("feed_id"))
	if err != nil {
		s.Logger.Info("failed to parse UUID from feed_id", "feed_id", feed_id)
		RespondWithError(w, 400, "invalid feed_id")
		return
	}

	f, _ := s.DB.FollowExists(r.Context(), database.FollowExistsParams{
		FeedID: feed_id,
		UserID: user_id,
	})
	log.Print(user_id.String(), feed_id.String(), f)
	err = s.DB.UnfollowFeed(r.Context(), database.UnfollowFeedParams{
		FeedID: feed_id,
		UserID: user_id,
	})
	if err != nil {
		s.Logger.Debug("failed to unfollow feed", "feed_id", feed_id, "user_id", user_id)
		RespondWithError(w, http.StatusInternalServerError, "failed to unfollow feed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *APIServer) GetUsersFeedFollows(w http.ResponseWriter, r *http.Request) {
	// get user_id
	dbUser, err := s.DB.GetUser(r.Context(), r.Context().Value(apiKeyContextKey).(string))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user := databaseUsertoUser(dbUser)

	dbFeedFollows, err := s.DB.GetUserFeeds(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload := databaseFeedFollowsToFeedFollows(dbFeedFollows)
	RespondWithJSON(w, http.StatusOK, payload)
}
