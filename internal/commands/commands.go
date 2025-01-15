package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"somewherecosmic/aggregator/internal/config"
	"somewherecosmic/aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

type State struct {
	Db   *database.Queries
	Conf *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	ValidCommands map[string]func(*State, Command) error
}

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

	s.Conf.Current_user = user.Name
	s.Conf.SetUser(user.Name)
	fmt.Printf("New user created: %s, Created_at: %s\n", user.Name, user.CreatedAt.Time)
	return nil
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
