package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/JoshuaTapp/gopherRSS/internal/database"
)

type APIServer struct {
	addr   string
	DB     *database.Queries
	Logger *slog.Logger
}

func NewAPIServer(addr, dbURL string) *APIServer {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	dbQueries := database.New(db)

	opts := &slog.HandlerOptions{
		// Use the ReplaceAttr function on the handler options
		// to be able to replace any single attribute in the log output
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// check that we are handling the time key
			if a.Key != slog.TimeKey {
				return a
			}

			t := a.Value.Time()

			// change the value from a time.Time to a String
			// where the string has the correct time format.
			a.Value = slog.StringValue(t.Format(time.DateTime))

			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return &APIServer{
		addr:   addr,
		DB:     dbQueries,
		Logger: logger,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	s.loadRoutes(router)

	server := http.Server{
		Addr:    s.addr,
		Handler: s.RequestLoggerMiddleware(router), // adding request middleware
	}

	// s.RunSetup(1000)

	log.Printf("Starting server at: %s", server.Addr)
	return server.ListenAndServe()
}

func (s *APIServer) loadRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /v1/healthz", ReadinessHandler)
	router.HandleFunc("GET /v1/err", ErrorHandler)

	router.HandleFunc("POST /v1/users", s.CreateUserHandler)
	router.Handle("GET /v1/users", s.RequireAuthnMiddleware(http.HandlerFunc(s.GetUserHandler)))

	router.Handle("POST /v1/feeds", s.RequireAuthnMiddleware(http.HandlerFunc(s.CreateFeedHandler)))
	router.HandleFunc("GET /v1/feeds", s.GetAllFeedsHandler)

	router.Handle("POST /v1/feed_follows", s.RequireAuthnMiddleware(http.HandlerFunc(s.FollowFeedHandler)))
	router.Handle("DELETE /v1/feed_follows/{feed_id}", s.RequireAuthnMiddleware(http.HandlerFunc(s.DeleteFeedFollowHandler)))
	router.Handle("GET /v1/feed_follows", s.RequireAuthnMiddleware(http.HandlerFunc(s.GetUsersFeedFollows)))

	router.Handle("GET /v1/posts/{limit}", s.RequireAuthnMiddleware(http.HandlerFunc(s.GetUsersPostsHandler)))
}
