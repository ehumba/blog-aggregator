package main

import (
	"fmt"

	"github.com/ehumba/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected 1 command input for login")
	}

	username := cmd.args[0]
	s.cfg.CurrentUserName = username

	fmt.Printf("username has been set to: %v", username)

	return nil
}
