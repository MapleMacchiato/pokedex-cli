package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MapleMacchiato/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Client, string) error
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
		"explore": {
			name:        "explore",
			description: "Get a list of pokemons from given area",
			callback:    commandExplore,
		},
	}
}

func commandExplore(pkc *pokeapi.Client, area string) error {
	fmt.Printf("Exploring %s\n", area)
	fmt.Printf("Found Pokemon:\n")
	return pkc.GetPokemon(area)
}

func commandMapB(pkc *pokeapi.Client, a string) error {
	if pkc.PrevURL == nil {
		return errors.New("No previous map")
	}
	return pkc.GetLocations(pkc.PrevURL)
}

func commandMap(pkc *pokeapi.Client, a string) error {
	return pkc.GetLocations(pkc.NextURL)
}

func commandHelp(pkc *pokeapi.Client, a string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Available list of Commands:")
	commands := getCommands()
	for k := range commands {
		fmt.Printf("name: %s\ndescription: %s\n\n", commands[k].name, commands[k].description)
	}
	return nil
}

func commandExit(pkc *pokeapi.Client, a string) error {
	fmt.Println("Exiting the Pokedex")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func run_repl(pkc *pokeapi.Client) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()
		words := cleanInput(reader.Text())
		commandName := words[0]
		var args string
		if len(words) > 1 {
			args = words[1]
		}
		command, ok := getCommands()[commandName]
		if !ok {
			fmt.Println("Invalid input, use help to see available commands")
		} else {
			err := command.callback(pkc, args)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	run_repl(&pokeClient)
}
