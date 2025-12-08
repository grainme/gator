package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
)

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
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
		return fmt.Errorf("couldn't create a user in the DB: %w\n", err)
	}

	// this edit this file ~/.gatorconfig.json
	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User has been created successfuly: %v\n", dbUser)
	return nil
}
