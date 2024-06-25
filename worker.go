package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const fetchInterval = time.Second * 5

func FetchFeedsWorker(s *APIServer, batchSize int) {
	s.Logger.Info("starting FetchFeedsWorker")
	tik := time.NewTicker(fetchInterval)

	for t := time.Now(); ; t = <-tik.C {
		fmt.Println("RUNNING WORKER at:", t)
		ctx := context.Background()
		// get next batchSize feeds
		dbFeeds, err := s.DB.GetNextFeedsToFetch(ctx, int32(batchSize))
		if err != nil {
			s.Logger.Warn("RSS Feed Worker", "error", err.Error())
			return
		}
		feeds := databaseFeedsToFeeds(dbFeeds)
		rss := s.FetchRssFeeds(feeds)

		// construct lookup table for reference
		feedUrlFeedIdMap := make(map[string]uuid.UUID)
		for _, feed := range feeds {
			feedUrlFeedIdMap[feed.Url] = feed.ID
		}

		// mark feeds fetched
		for _, f := range rss {
			err := s.DB.UpdateFeedFetchTime(ctx, f.url)
			if err != nil {
				s.Logger.Warn("failed to mark feed as fetched", "URL", f.url, "error", err.Error())
				continue
			}
			for _, item := range f.rss.Channel.Items {
				err = s.CreatePost(item, feedUrlFeedIdMap[f.url])
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}
