package commands

import (
	"errors"
	"fmt"
	"somewherecosmic/aggregator/internal/config"
)

type State struct {
	Conf *config.Config
}

type Command struct {
	Name string
	args []string
}

type Commands struct {
	validCommands map[string]func(*State, Command) error
}

func handerLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("function expects a single argument: username")
	}

	if err := s.Conf.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Printf("User has been set: %s\n", cmd.args[0])

	return nil
}

func (c *Commands) register(name string, f func(*State, Command) error) error {
	c.validCommands[name] = f

	return nil
}

func (c *Commands) run(s *State, cmd Command) error {
	err := c.validCommands[cmd.Name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}
