package cli

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: login <name>")
	}

	username := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist in the db: %w", err)
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	slog.Info("user logged in", "name", user.Name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: register <name>")
	}

	name := cmd.Args[0]
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	_, err := s.Db.GetUser(context.Background(), name)
	if err == nil {
		// then user already exist in the db
		return fmt.Errorf("user already exists in the DB: %w\n", err)
	}

	dbUser, err := s.Db.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to create user %q: %w", name, err)
	}

	// this edit this file ~/.gatorconfig.json
	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	slog.Info("user created", "name", dbUser.Name, "id", dbUser.ID)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: reset")
	}

	rowsDeleted, err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users from the db: %v", err)
	}

	slog.Info("users reset", "rows_deleted", rowsDeleted)
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: users")
	}

	dbUsers, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve users from the database: %v", err)
	}

	for _, user := range dbUsers {
		isCurrent := s.Cfg.CurrentUserName == user.Name
		slog.Info("user", "name", user.Name, "current", isCurrent)
	}

	return nil
}
