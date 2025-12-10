package internal

import (
	"context"
	"fmt"

	"github.com/grainme/gator/internal/database"
)

func HandlerUnfollow(s *State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: unfollow <feed_url>\n")
	}

	feedURL := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return err
	}

	// delete a feed follow of current user
	_, err = s.Db.DeleteByUserIdAndFeedId(context.Background(), database.DeleteByUserIdAndFeedIdParams{
		UserID: currentUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
