package cli

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/grainme/gator/internal/database"
)

const (
	DefaultPostLimit = 2
)

func HandlerBrowse(s *State, cmd Command, currentUser database.User) error {
	limit := DefaultPostLimit

	if len(cmd.Args) == 1 {
		if userLimit, err := strconv.Atoi(cmd.Args[0]); err == nil && userLimit > 0 {
			limit = userLimit
		} else {
			return err
		}
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
		slog.Info("no posts available")
		return nil
	}
	for _, post := range posts {
		slog.Info("post", "title", post.Title, "url", post.Url, "description", post.Description, "published_at", post.PublishedAt)
	}
	return nil
}
