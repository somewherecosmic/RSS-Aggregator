package rss

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

func (feed *RSSFeed) Stringify() {
	fmt.Printf("Title: %s\nLink: %s\nDescription: %s\nItems: %s\n", feed.Channel.Title, feed.Channel.Link, feed.Channel.Description, feed.Channel.Item)
}

// func (feed *RSSFeed) decodeEscapedChars() {
// 	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
// 	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
// 	for _, item := range feed.Channel.Item {
// 		item.Title = html.UnescapeString(item.Title)
// 		item.Description = html.UnescapeString(item.Description)
// 	}
// }

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (item *RSSItem) Stringify() {
	fmt.Printf("Title: %s\n Link: %s\n, Description: %s\n, PubDate: %s\n", item.Title, item.Link, item.Description, item.PubDate)
}

// Fetch feed from url
// return a RSSFeed struct

// Create new request with http client
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	feed := RSSFeed{}
	if err := xml.Unmarshal(bytes, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}
	return &feed, nil
}
