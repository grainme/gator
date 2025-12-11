package internal

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grainme/gator/internal/database"
)

func HandlerBrowse(s *State, cmd Command, currentUser database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		arg, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = arg
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <limit>", cmd.Name)
	}

	posts, err := s.Db.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: currentUser.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		fmt.Println("No posts available!")
		return nil
	}
	for _, post := range posts {
		fmt.Println(post.Title)
	}
	return nil
}
