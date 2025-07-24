package main

import (
	"context"
	"fmt"

	"github.com/ehumba/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		if s.cfg.CurrentUserName == "" {
			fmt.Println("no user logged in")
			return nil
		}

		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, c, currentUser)
	}
}
