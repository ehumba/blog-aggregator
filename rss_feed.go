package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// new request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("could not create new request: %v", err)
	}

	// set header
	req.Header.Set("User-Agent", "gator")

	// make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	// read from the response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// unmarshal data into RSSFeed struct
	var feed RSSFeed
	if err = xml.Unmarshal(data, &feed); err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return &feed, nil
}

func unescape(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}
