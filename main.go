package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/grainme/gator/internal"
	"github.com/grainme/gator/internal/config"
	"github.com/grainme/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	commands := internal.Commands{
		RegistredCommands: map[string]func(*internal.State, internal.Command) error{},
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	state := internal.State{
		Cfg: cfg,
		Db:  dbQueries,
	}

	if err := commands.Register("login", internal.HandlerLogin); err != nil {
		log.Fatalf("error registering login command: %v", err)
	}
	if err := commands.Register("register", internal.HandlerRegister); err != nil {
		log.Fatalf("error registering register command: %v", err)
	}
	if err := commands.Register("reset", internal.HandlerReset); err != nil {
		log.Fatalf("error registering reset command: %v", err)
	}
	if err := commands.Register("users", internal.HandlerGetUsers); err != nil {
		log.Fatalf("error registering users command: %v", err)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: command name is missing")
		os.Exit(1)
	}

	cmd := internal.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := commands.Run(&state, cmd); err != nil {
		log.Fatalf("error running command %q: %v", cmd.Name, err)
	}
}
