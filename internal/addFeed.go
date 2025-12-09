package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
)

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>\n")
	}

	// parsing args
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	currentUser, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not fetch current user: %v", err)
	}

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

	fmt.Println(feedCreated)
	return nil
}
