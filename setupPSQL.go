package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/JoshuaTapp/gopherRSS/internal/database"
	"github.com/google/uuid"
)

func (s *APIServer) RunSetup(n int) error {
	userIds := s.PopulateUsers(n)
	feedIds := s.PopulateFeeds(userIds)
	userFollows := SetupFeedFollows(userIds, feedIds)
	err := s.PopulateFeedFollows(userFollows)

	return err
}

func (s *APIServer) PopulateUsers(n int) (userIds []uuid.UUID) {
	var userParams []database.CreateUserParams

	for i := 1; i <= n; i++ {
		userParams = append(userParams, database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      fmt.Sprintf("User %d", i),
		})
	}

	for _, p := range userParams {
		user, err := s.DB.CreateUser(context.Background(), p)
		if err != nil {
			fmt.Println("Error: PopulateUsers() failed DB call - ", err.Error())
			return
		}
		userIds = append(userIds, user.ID)
	}

	fmt.Println("Success: PopulateUsers() created users")
	for _, id := range userIds {
		fmt.Printf("\tuserId: %s\n", id.String())
	}
	fmt.Println()
	return
}

func (s *APIServer) PopulateFeeds(userIds []uuid.UUID) (feedIds []uuid.UUID) {
	feedParams := []struct {
		Url  string
		Name string
	}{
		{
			"https://blog.boot.dev/index.xml",
			"Boot.Dev",
		},
		{
			"http://rss.cnn.com/rss/cnn_topstories.rss",
			"CNN Top Stories",
		},
		{
			"https://wagslane.dev/index.xml",
			"Lane's Blog",
		},
		{
			"https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml",
			"NYT Tech Stories",
		},
		{
			"https://feeds.bbci.co.uk/news/world/us_and_canada/rss.xml?edition=int",
			"BBC America",
		},
		{
			"https://feeds.feedburner.com/animals",
			"Animal of the Day!",
		},
	}

	dbFeedParams := make([]database.CreateFeedParams, 0)
	for _, fp := range feedParams {
		dbFeedParams = append(dbFeedParams, database.CreateFeedParams{
			ID:        uuid.New(),
			Name:      fp.Name,
			Url:       fp.Url,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    userIds[rand.Intn(len(userIds))],
		})
	}

	// create feeds
	for i := range dbFeedParams {
		fmt.Printf("Adding Feed: %v\n%v\n", dbFeedParams[i].Name, dbFeedParams[i])
		dbFeed, err := s.DB.CreateFeed(context.Background(), dbFeedParams[i])
		if err != nil {
			fmt.Println("Error: PopulateFeeds() failed DB call - ", err.Error())
			return
		}
		feedIds = append(feedIds, databaseFeedToFeed(dbFeed).ID)
	}

	return
}

type userFeedID struct {
	UserId uuid.UUID
	FeedId uuid.UUID
}

func SetupFeedFollows(userIds, feedIds []uuid.UUID) (pairs []userFeedID) {
	if len(feedIds) < 1 {
		fmt.Println("feedIds < 1 length")
		return
	}

	for _, id := range userIds {
		followCnt := rand.Intn(len(feedIds))

		sample := func() []uuid.UUID {
			rand.Shuffle(len(feedIds), func(i, j int) {
				feedIds[i], feedIds[j] = feedIds[j], feedIds[i]
			})
			return feedIds[:followCnt]
		}()

		for _, s := range sample {
			pairs = append(pairs, userFeedID{UserId: id, FeedId: s})
		}

	}

	return
}

func (s *APIServer) PopulateFeedFollows(feedFollows []userFeedID) error {
	for _, follow := range feedFollows {
		_, err := s.DB.FollowFeed(context.Background(), database.FollowFeedParams{
			FeedID:    follow.FeedId,
			UserID:    follow.UserId,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			fmt.Println("Error: PopulateFeedFollows() failed DB call - ", err.Error())
			return err
		}
	}
	fmt.Println("Success: PopulateFeedFollows() added feed follows")
	return nil
}
