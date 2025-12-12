package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/grainme/gator/internal/cli"
	"github.com/grainme/gator/internal/config"
	"github.com/grainme/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	commands := cli.Commands{
		RegistredCommands: map[string]func(*cli.State, cli.Command) error{},
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	state := cli.State{
		Cfg: cfg,
		Db:  dbQueries,
	}

	if err := commands.Register("login", cli.HandlerLogin); err != nil {
		log.Fatalf("error registering login command: %v", err)
	}
	if err := commands.Register("register", cli.HandlerRegister); err != nil {
		log.Fatalf("error registering register command: %v", err)
	}
	if err := commands.Register("reset", cli.HandlerReset); err != nil {
		log.Fatalf("error registering reset command: %v", err)
	}
	if err := commands.Register("users", cli.HandlerGetUsers); err != nil {
		log.Fatalf("error registering users command: %v", err)
	}
	if err := commands.Register("agg", cli.HandlerAggregator); err != nil {
		log.Fatalf("error registering agg command: %v", err)
	}
	if err := commands.Register("addfeed", cli.MiddlewareLoggedIn(cli.HandlerAddFeed)); err != nil {
		log.Fatalf("error registering addfeed command: %v", err)
	}
	if err := commands.Register("feeds", cli.HandlerGetFeeds); err != nil {
		log.Fatalf("error registering feeds command: %v", err)
	}
	if err := commands.Register("follow", cli.MiddlewareLoggedIn(cli.HandlerFollow)); err != nil {
		log.Fatalf("error registering follow command: %v", err)
	}
	if err := commands.Register("following", cli.MiddlewareLoggedIn(cli.HandlerFollowing)); err != nil {
		log.Fatalf("error registering following command: %v", err)
	}
	if err := commands.Register("unfollow", cli.MiddlewareLoggedIn(cli.HandlerUnfollow)); err != nil {
		log.Fatalf("error registering unfollow command: %v", err)
	}
	if err := commands.Register("browse", cli.MiddlewareLoggedIn(cli.HandlerBrowse)); err != nil {
		log.Fatalf("error registering browse command: %v", err)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: command name is missing")
		os.Exit(1)
	}

	cmd := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := commands.Run(&state, cmd); err != nil {
		log.Fatalf("error running command %q: %v", cmd.Name, err)
	}
}
