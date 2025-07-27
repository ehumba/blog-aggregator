package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ehumba/blog-aggregator/internal/database"
)

func scrapeFeeds(s *state) {
	// get next feed to fetch and mark it as fetched
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("unable to get next feed to fetch: %v", err)
		return
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now()},
		UpdatedAt:     time.Now(),
		ID:            nextFeed.ID})
	if err != nil {
		fmt.Print(err)
		return
	}

	// fetch the feed
	feedOutput, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Print(err)
		return
	}

	unescape(feedOutput)

	// print content
	fmt.Println(feedOutput.Channel.Title)
	fmt.Println(feedOutput.Channel.Description)
	for _, item := range feedOutput.Channel.Item {
		fmt.Println(item.Title)
		fmt.Println(item.Description)
	}

}
