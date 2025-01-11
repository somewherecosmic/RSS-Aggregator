package commands

import (
	"fmt"
	"log"
	"somewherecosmic/aggregator/internal/config"
)

type State struct {
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

	if err := s.Conf.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("User has been set: %s\n", cmd.Args[0])

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
