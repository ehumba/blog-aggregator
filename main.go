package main

import (
	"fmt"
	"os"

	"github.com/ehumba/blog-aggregator/internal/config"
)

func main() {
	// read config file
	conf, err := config.Read()
	if err != nil {
		panic(err)
	}

	currentState := state{
		cfg: &conf,
	}

	//Initialize the commands map
	cmds := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	// Get command-line arguments passed by the user
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	inputCmdName := argsWithProg[1]
	inputArgs := argsWithProg[2:]

	cmds.run(&currentState, command{name: inputCmdName, args: inputArgs})

}
