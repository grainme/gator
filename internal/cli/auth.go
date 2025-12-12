package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grainme/gator/internal/database"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist in the db: %w", err)
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User has been set successfuly")
	return nil
}

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

func HandlerReset(s *State, _ Command) error {
	rowsDeleted, err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users from the db: %v", err)
	}

	fmt.Printf("Successfuly deleted %d rows from USERS table\n", rowsDeleted)
	return nil
}

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
