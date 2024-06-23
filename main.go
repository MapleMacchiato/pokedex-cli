package main

import (
	"bufio"
	"fmt"
	"github.com/MapleMacchiato/pokedex-cli/internal/pokeapi"
	"os"
	"strings"
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
	fmt.Println("Exit")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func run_repl() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()
		words := cleanInput(reader.Text())
		commandName := words[0]
		command, ok := getCommands()[commandName]
		if !ok {
			fmt.Println("Invalid input, use help to see available commands")
		} else {
			command.callback()
		}
	}
}

func main() {
	run_repl()
}
