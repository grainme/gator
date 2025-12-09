package internal

import (
	"context"
	"fmt"
)

func HandlerGetUsers(s *State, _ Command) error {
	dbUsers, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve users from the database: %v", err)
	}

	for _, user := range dbUsers {
		if s.Cfg.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}
