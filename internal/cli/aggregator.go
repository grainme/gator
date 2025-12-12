package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
	"github.com/grainme/gator/internal/rss"
)

func HandlerAggregator(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err := ScrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func ScrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("----------------------")
	fmt.Println("Found a feed to fetch!")
	err = s.Db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return nil
	}

	feedItems, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	// convert pubDate (string) to time
	// this format: Mon, 01 Jan 0001 00:00:00 +0000
	for _, item := range feedItems.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}
		post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}
		_, err = s.Db.CreatePost(context.Background(), post)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate") {
				return err
			}
		}
	}

	return nil
}
