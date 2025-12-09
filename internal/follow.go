package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
)

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	username := s.Cfg.CurrentUserName
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	feedURL := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	createdFeedFollow, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println(createdFeedFollow.Feedname, createdFeedFollow.Username)
	return nil
}
