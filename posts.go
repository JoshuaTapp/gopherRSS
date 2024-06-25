package main

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (s *APIServer) CreatePost(i Item, feed_id uuid.UUID) error {
	dupErr := errors.New("pq: duplicate key value violates unique constraint \"posts_url_key\"").Error()

	post := Post{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       &i.Title,
		Url:         i.Link,
		Description: &i.Desc,
		PublishedAt: ParsePubDate(i.PubDate),
		FeedID:      feed_id,
	}

	_, err := s.DB.CreatePost(context.Background(), postToDatabasePostParams(post))
	if err != nil && err.Error() != dupErr {
		s.Logger.Warn("Error creating post", "url", i.Link, "feedId", feed_id)
		return err
	}
	return nil
}
