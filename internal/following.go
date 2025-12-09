package internal

import (
	"context"
	"fmt"
)

func HandlerFollowing(s *State, _ Command) error {
	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
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
