package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ehumba/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

	// get published time
	for _, item := range feedOutput.Channel.Item {
		var publishedAt time.Time
		var parseErr error

		layouts := []string{
			time.RFC1123Z,
			time.RFC1123,
			time.RFC822Z,
			time.RFC822,
			time.RFC3339,
		}

		for _, layout := range layouts {
			publishedAt, parseErr = time.Parse(layout, item.PubDate)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			fmt.Printf("could not parse published time: %v", parseErr)
			publishedAt = time.Now()
		}

		// create the post
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}
			fmt.Printf("failed to create post: %v\n", err)
		}
	}
}

// helper function to check if the error is a unique violation
func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
