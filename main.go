package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ehumba/blog-aggregator/internal/config"
	"github.com/ehumba/blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// read config file
	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", "postgres://postgres:postgress@localhost:5432/gator")
	if err != nil {
		fmt.Printf("could not open database: %v", err)
	}
	dbQueries := database.New(db)

	currentState := state{
		db:  dbQueries,
		cfg: &conf,
	}

	//Initialize the commands map
	cmds := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	// Get command-line arguments passed by the user
	argsWithProg := os.Args

	if len(argsWithProg) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	inputCmdName := argsWithProg[1]
	inputArgs := argsWithProg[2:]

	err = cmds.run(&currentState, command{name: inputCmdName, args: inputArgs})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
