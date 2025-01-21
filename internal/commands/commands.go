package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"somewherecosmic/aggregator/internal/config"
	"somewherecosmic/aggregator/internal/database"
	"somewherecosmic/aggregator/internal/rss"
	"time"

	"github.com/google/uuid"
)

/*
Structs + Struct methods - State, Command, Commands
*/
type State struct {
	Db   *database.Queries
	Conf *config.Config
}

// func (s *State) setUser(user *database.User) error {
// 	if user == nil {
// 		return errors.New("user object passed is nil")
// 	}
// 	s.User = user
// 	return nil
// }

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	ValidCommands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	c.ValidCommands[name] = f
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	err := c.ValidCommands[cmd.Name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}

/*
	Command handler functions
*/

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		log.Fatal("function expects a single argument: username")
	}

	_, err := s.Db.FindUserByName(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("User doesn't exist")
		os.Exit(1)
	}

	if err := s.Conf.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("User has been set: %s\n", cmd.Args[0])

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		log.Fatal("Command expects username as an argument")
	}

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: cmd.Args[0],
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.Conf.SetUser(user.Name)
	fmt.Printf("New user created: %s, Created_at: %s\n", user.Name, user.CreatedAt.Time)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) > 1 {
		return errors.New("too many arguments given - reset expects 0 arguments")
	}

	if err := s.Db.ClearUserTable(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("User table reset - all records deleted")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Args) > 1 {
		return errors.New("too many arguments given, users expects 0")
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		current := ""
		if user.Name == s.Conf.Current_user {
			current = "(current)"
		}
		fmt.Printf("* %s "+current+"\n", user.Name)
	}

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("addfeed expects two arguments: name, url")
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: cmd.Args[0],
		Url:  cmd.Args[1],
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	HandlerFollow(s, Command{"follow", []string{feed.Url}}, user)

	return nil
}

func HandlerFeed(s *State, cmd Command) error {
	if len(cmd.Args) > 1 {
		return errors.New("too many arguments provided - agg expects 0")
	}

	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	// feed.Stringify()
	fmt.Println(feed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Created By: %s\n", feed.User)
	}

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("%s takes no arguments", cmd.Name)
	}

	followedFeeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range followedFeeds {
		fmt.Println(feed.FeedName)
	}

	return nil
}

// takes a single url
// creates a new feed follow record for a user

func HandlerFollow(s *State, cmd Command, user database.User) error {
	// fmt.Println("Before user access")
	// fmt.Printf("User ID: %v\n", s.User.Name)
	// fmt.Println("After user access")
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Feed ID: %s\n", feed.ID)

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("%s has just followed the feed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	if err := s.Db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	}); err != nil {
		return err
	}

	return nil
}

// type HandlerFunc func(s *State, cmd Command) error

// type HandlerFuncLoggedIn func(s *State, cmd Command, user database.User) error

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {

	return func(s *State, cmd Command) error {
		user, err := s.Db.FindUserByName(context.Background(), s.Conf.Current_user)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
