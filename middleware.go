package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

// Define a new type for context keys
type contextKey string

// Create a key with the new type
const apiKeyContextKey contextKey = "apiKey"

func getApiKey(head *http.Header) (string, error) {
	apikeyHeader := head.Get("Authorization")
	if apikeyHeader == "" {
		return "", errors.New("apikey not found")
	}
	splitHead := strings.Split(apikeyHeader, " ")
	if len(splitHead) < 2 || splitHead[0] != "ApiKey" {
		return "", errors.New("mailformed authorization header")
	}

	return splitHead[1], nil
}

func (s *APIServer) RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Before Logic
		key, err := getApiKey(&r.Header)
		if err != nil {
			s.Logger.InfoContext(
				r.Context(),
				"request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
		} else {
			s.Logger.InfoContext(
				r.Context(),
				"request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("apiKey", key),
			)
		}
		// continue request processing
		next.ServeHTTP(w, r)

		// After
		
	}
}

func (s *APIServer) RequireAuthnMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if user is authenticated
		key, err := getApiKey(&r.Header)
		if err != nil {
			// failed auth
			RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		// auth exist in database?
		isUser, err := s.DB.IsUser(r.Context(), key)
		if err != nil || !isUser {
			// failed auth
			RespondWithError(w, http.StatusForbidden, "invalid user")
			return
		}

		ctx := context.WithValue(r.Context(), apiKeyContextKey, key)
		// valid user, add apiKey to context and process
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
