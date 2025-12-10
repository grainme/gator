package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
)

func HandlerAddFeed(s *State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <feed_name> <feed_url>\n")
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUser.ID,
	}
	feedCreated, err := s.Db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	// create a feed follow (current user following that feed)
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feedCreated.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println(feedCreated)
	return nil
}
