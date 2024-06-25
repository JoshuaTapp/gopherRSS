// feed_test.go
package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockAPIServer struct {
	APIServer
}

// Helper function to create feeds
func createFeeds() []Feed {
	return []Feed{
		{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "Boot.Dev",
			Url:           "https://blog.boot.dev/index.xml",
			UserID:        uuid.New(),
			LastFetchedAt: nullTimeToTimePtr(sql.NullTime{}),
		},
		{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "CNN Top Stories",
			Url:           "http://rss.cnn.com/rss/cnn_topstories.rss",
			UserID:        uuid.New(),
			LastFetchedAt: nullTimeToTimePtr(sql.NullTime{}),
		},
		{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "Lane's Blog",
			Url:           "https://wagslane.dev/index.xml",
			UserID:        uuid.New(),
			LastFetchedAt: nullTimeToTimePtr(sql.NullTime{}),
		},
		{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "NYT Tech Stories",
			Url:           "https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml",
			UserID:        uuid.New(),
			LastFetchedAt: nullTimeToTimePtr(sql.NullTime{}),
		},
		{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "BBC America",
			Url:           "https://feeds.bbci.co.uk/news/world/us_and_canada/rss.xml?edition=int",
			UserID:        uuid.New(),
			LastFetchedAt: nullTimeToTimePtr(sql.NullTime{}),
		},
		{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "Animal of the Day!",
			Url:           "https://feeds.feedburner.com/animals",
			UserID:        uuid.New(),
			LastFetchedAt: nullTimeToTimePtr(sql.NullTime{}),
		},
	}
}

// Mock HTTP Response for RSS feed
func mockHTTPResponse(feedURL string, response string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		if r.URL.String() == feedURL {
			w.Write([]byte(response))
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})
	return httptest.NewServer(handler)
}

func TestFetchRssFeeds(t *testing.T) {
	s := &mockAPIServer{
		APIServer: APIServer{
			Logger: slog.Default(),
		},
	}

	// Example mock server for a specific URL
	testFeedURL := "/rss/testfeed"
	testFeedResponse := `<?xml version="1.0" encoding="UTF-8"?><rss><channel><title>Test Feed</title></channel></rss>`
	mockServer := mockHTTPResponse(testFeedURL, testFeedResponse)
	defer mockServer.Close()

	// Adjust the URL in feeds to point to the mock server
	feeds := createFeeds()
	// For real world test: comment out the below for range loop
	for i := range feeds {
		feeds[i].Url = mockServer.URL + testFeedURL
	}

	testCases := []struct {
		name    string
		fetchFn func([]Feed) []FeedUrlRSS
	}{
		{"FetchRssFeeds (goroutines)", s.FetchRssFeeds},
		{"FetchRssFeedsSlow (no goroutines)", s.FetchRssFeedsSlow},
	}

	const repetitions = 10
	results := map[string][]time.Duration{
		"FetchRssFeeds (goroutines)":        {},
		"FetchRssFeedsSlow (no goroutines)": {},
	}

	for _, tc := range testCases {
		for i := 0; i < repetitions; i++ {
			t.Run(fmt.Sprintf("%s - run %d", tc.name, i), func(t *testing.T) {
				start := time.Now()
				tc.fetchFn(feeds)
				elapsed := time.Since(start)
				results[tc.name] = append(results[tc.name], elapsed)
			})
		}
	}

	// Calculate and print averages
	for name, timings := range results {
		var totalDuration time.Duration
		for _, duration := range timings {
			totalDuration += duration
		}
		averageDuration := totalDuration / time.Duration(repetitions)
		fmt.Printf("%s - Average Time taken: %s\n", name, averageDuration.String())
	}
}
