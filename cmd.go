package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ehumba/blog-aggregator/internal/config"
	"github.com/ehumba/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type state struct {
	db  *database.Queries
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

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist: %v", err)
	}

	s.cfg.SetUser(username)

	fmt.Printf("username has been set to: %v\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no name provided")
	}
	username := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		// check for unique violation first
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("user with name %s already exists", username)
			}
		}
		return fmt.Errorf("could not register user: %v", err)
	}

	s.cfg.SetUser(username)
	fmt.Printf("New user %v successfully registered\n", username)
	fmt.Printf("%+v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ClearTable(context.Background())
	if err != nil {
		return fmt.Errorf("database reset failed: %v", err)
	}
	fmt.Println("database reset successful")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failure: %v", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	// parse the argument into a Duration
	if len(cmd.args) != 1 {
		return fmt.Errorf("one argument required: <duration> (example: 10s)")
	}
	timeBetweenReqs := cmd.args[0]
	parsedTime, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("invalid duration string. Please enter a duration in the format 1s, 1m, 1h, etc")
	}

	// start the loop of scraping feeds with a ticker
	fmt.Printf("Collecting feeds every %v\n", parsedTime)
	ticker := time.NewTicker(parsedTime)
	defer ticker.Stop()
	scrapeFeeds(s)
	for range ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("two arguments required: <name> <url>")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("falied to create new Feed: %v", err)
	}
	fmt.Printf("New feed added!\nID: %v\nCreated at: %v\nUpdated at: %v\nName: %v\nURL: %v\nUser ID: %v\n", newFeed.ID, newFeed.CreatedAt, newFeed.UpdatedAt, newFeed.Name, newFeed.Url, newFeed.UserID)

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow: %v", err)
	}
	fmt.Printf("%s successfully followed %s!\n", s.cfg.CurrentUserName, newFeed.Name)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %v\nUrl: %v\n", feed.Name, feed.Url)
		for _, user := range users {
			if user.ID == feed.UserID {
				fmt.Printf("Username: %v\n\n", user.Name)
			}
		}
	}
	return nil
}

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("one argument required: <url>")
	}

	feedToFollow, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("feed not found")
	}

	newFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feedToFollow.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow: %v", err)
	}

	fmt.Printf("%s successfully followed %s\n", newFollow.UserName, newFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, currentUser database.User) error {
	following, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s follows:\n", s.cfg.CurrentUserName)

	for _, follow := range following {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("one argument required: <url>")
	}
	url := cmd.args[0]

	err := s.db.DeleteFollow(context.Background(), database.DeleteFollowParams{
		Url:  url,
		Name: currentUser.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow: %v", url)
	}
	fmt.Printf("%s successfully unfollowed %s\n", s.cfg.CurrentUserName, url)
	return nil
}

func handlerBrowse(s *state, cmd command, currentUser database.User) error {
	limit := "2"
	if len(cmd.args) > 0 {
		limit = cmd.args[0]
	}

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		return fmt.Errorf("invalid limit: %v", err)
	}

	fmt.Println("You may enter a limit for displayed posts. Default: 2")

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: currentUser.ID,
		Limit:  int32(parsedLimit),
	})
	if err != nil {
		return fmt.Errorf("failed to get posts: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %v\n", post.Name)
		fmt.Printf("Post URL: %v\n", post.Url)
		fmt.Printf("Published at: %v\n", post.PublishedAt)
		if post.Description.Valid {
			fmt.Printf("Description: %v\n", post.Description.String)
		}
		fmt.Println("---------------------")
	}

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %v", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
