package internal

import (
	"context"
	"fmt"
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


