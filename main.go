package main

import (
	"fmt"
	"log"
	"os"
	"somewherecosmic/aggregator/internal/commands"
	"somewherecosmic/aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	var state commands.State
	state.Conf = conf

	commandRegistry := commands.Commands{
		ValidCommands: make(map[string]func(*commands.State, commands.Command) error),
	}

	commandRegistry.Register("login", commands.HandlerLogin)

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatal("Error: Too few arguments provided")
	}

	issuedCommand := commands.Command{
		Name: userArgs[1],
		Args: userArgs[2:],
	}

	if err := commandRegistry.Run(&state, issuedCommand); err != nil {
		fmt.Println(err)
	}
}
