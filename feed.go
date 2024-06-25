package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"
)

type FeedUrlRSS struct {
	url string
	rss Rss
}

func (s *APIServer) FetchRssFeed(url string) (data FeedUrlRSS, err error) {
	r, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer r.Body.Close()

	rss := Rss{}
	decoder := xml.NewDecoder(r.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return data, err
	}
	return FeedUrlRSS{url, rss}, nil
}

func (s *APIServer) FetchRssFeeds(feeds []Feed) []FeedUrlRSS {
	var wg sync.WaitGroup

	rssCh := make(chan FeedUrlRSS, len(feeds))
	errsCh := make(chan error, len(feeds))

	for _, feed := range feeds {
		wg.Add(1)

		go func(f Feed) {
			defer wg.Done()
			data, err := s.FetchRssFeed(f.Url)
			if err != nil {
				errsCh <- err
				return
			}

			rssCh <- data
		}(feed)
	}

	wg.Wait()
	close(errsCh)
	close(rssCh)

	rssFeeds := make([]FeedUrlRSS, 0)

	for r := range rssCh {
		fmt.Printf("RSS Feed - Title: %s, URL: %s\n", r.rss.Channel.Title, r.url)
		rssFeeds = append(rssFeeds, r)
	}

	for e := range errsCh {
		s.Logger.Warn("RSS error", "err", e.Error())
	}

	return rssFeeds
}

func (s *APIServer) FetchRssFeedsSlow(feeds []Feed) []FeedUrlRSS {
	rssFeeds := make([]FeedUrlRSS, 0)

	for _, f := range feeds {
		r, err := s.FetchRssFeed(f.Url)
		if err != nil {
			fmt.Printf("ERROR!\n\tURL: %s\n\terr: %v\n", f.Url, err.Error())
			continue
		}
		rssFeeds = append(rssFeeds, r)
		fmt.Printf("RSS Feed - Title: %s, URL: %s\n", r.rss.Channel.Title, r.url)
	}

	return rssFeeds
}
