package main

import (
	"net/http"
	"strconv"

	"github.com/JoshuaTapp/gopherRSS/internal/database"
)

const (
	minLimit = 1
	maxLimit = 1024
)

func (s *APIServer) GetUsersPostsHandler(w http.ResponseWriter, r *http.Request) {
	// get limit from query params
	// get user_id
	user, err := s.DB.GetUser(r.Context(), r.Context().Value(apiKeyContextKey).(string))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	limit := func() int32 {
		v := r.PathValue("limit")
		l, err := strconv.Atoi(v)
		if err != nil {
			return minLimit
		}
		if l < 1 {
			return minLimit
		}
		if l > maxLimit {
			return maxLimit
		}
		return int32(l)
	}

	// get user's posts
	dbPosts, err := s.DB.GetUsersPosts(r.Context(), database.GetUsersPostsParams{
		UserID: user.ID,
		Limit:  limit(),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	posts := databasePostsToPosts(dbPosts)
	RespondWithJSON(w, 200, posts)
}
