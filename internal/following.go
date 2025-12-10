package internal

import (
	"context"
	"fmt"

	"github.com/grainme/gator/internal/database"
)

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
