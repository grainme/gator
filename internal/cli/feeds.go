package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
)

func HandlerGetFeeds(s *State, _ Command) error {
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(feed.Name, feed.Url, user.Name)
	}

	return nil
}

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

func HandlerFollow(s *State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
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
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}
	createdFeedFollow, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println(createdFeedFollow.Feedname, createdFeedFollow.Username)
	return nil
}

func HandlerFollowing(s *State, _ Command, currentUser database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollows {
		feedFollowExtraInfos, err := s.Db.GetFeedFollowByFeedId(context.Background(), feedFollow.FeedID)
		if err != nil {
			return err
		}
		fmt.Println(feedFollowExtraInfos.Feedname)
	}

	return nil
}

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
