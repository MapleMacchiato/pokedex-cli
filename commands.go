package main

import (
	"fmt"
	"github.com/MapleMacchiato/pokedex-cli/internal/pokeapi"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Get next locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get previous locations",
			callback:    commandMapB,
		},
	}
}

func commandMapB() error {
	pokeapi.GetMapsB()
	return nil
}

func commandMap() error {
	pokeapi.GetMaps()
	return nil
}
func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available list of Commands:")
	commands := getCommands()
	for k := range commands {
		fmt.Printf("name: %s\ndescription: %s\n\n", commands[k].name, commands[k].description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Exiting the Pokedex")
	os.Exit(0)
	return nil
}
