package internal

import (
	"fmt"

	"github.com/grainme/gator/internal/config"
	"github.com/grainme/gator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db  *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegistredCommands map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	// check if cmd is supported
	handler, exists := c.RegistredCommands[cmd.Name]
	if !exists {
		return fmt.Errorf("command (%v) is not supported\n", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *Commands) Register(cmd string, f func(*State, Command) error) error {
	if cmd == "" {
		return fmt.Errorf("command name should not be empty\n")
	}
	c.RegistredCommands[cmd] = f
	return nil
}
