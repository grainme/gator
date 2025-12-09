package internal

import (
	"context"
	"fmt"
	"html"
)

func HandlerAggregator(s *State, _ Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	// decoding escaped HTML entities (like &ldquo)
	feed.Channel.Descriptin = html.UnescapeString(feed.Channel.Descriptin)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	fmt.Printf("%+v\n", feed)

	return nil
}
