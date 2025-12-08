package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/grainme/gator/internal"
	"github.com/grainme/gator/internal/config"
	"github.com/grainme/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error: reading file - \n%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	// init(ing) structures
	state := internal.State{
		Cfg: cfg,
	}
	commands := internal.Commands{
		RegistredCommands: map[string]func(*internal.State, internal.Command) error{},
	}

	// setting up the DB
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Printf("error: couldn't connect to postgres (more info below) \n-- %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries

	// support "login" command
	err = commands.Register("login", internal.HandlerLogin)
	if err != nil {
		fmt.Printf("error: saving command (more info below) \n-- %v\n", err)
		os.Exit(1)
	}
	// support "register" command
	err = commands.Register("register", internal.HandlerRegister)
	if err != nil {
		fmt.Printf("error: saving command (more info below) \n-- %v\n", err)
		os.Exit(1)
	}

	// os.Args[0] is just the cli name - skipi :)
	var commandName string
	var args []string = make([]string, 0)
	// command name is not optional (at least for now)
	if len(os.Args) < 2 {
		fmt.Printf("error: command name is missing\n")
		os.Exit(1)
	}
	commandName = os.Args[1]
	// args are optional for some commands
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	cmd := internal.Command{
		Name: commandName,
		Args: args,
	}

	err = commands.Run(&state, cmd)
	if err != nil {
		fmt.Printf("error: running command (more info below) \n-- %v\n", err)
		os.Exit(1)
	}

}
