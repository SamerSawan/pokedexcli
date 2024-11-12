package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/samersawan/pokedexcli/internal/api"
)

type config struct {
	next string
	prev *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
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
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent call to map displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func commandHelp(cfg *config) error {
	commands := getCommands()
	commandOrder := []string{"help", "exit", "map", "mapb"}
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println()
	for _, name := range commandOrder {
		fmt.Println(commands[name].name + ": " + commands[name].description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *config) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {
	prev, next, locations, err := api.GetLocations(cfg.next)
	if err != nil {
		return err
	}
	cfg.next = next
	cfg.prev = prev
	for i := 0; i < len(locations); i++ {
		fmt.Println(locations[i])
	}
	return nil
}

func commandMapb(cfg *config) error {

	if cfg.prev == nil {
		fmt.Println("Can not display previous locations because they do not exist. Use map instead.")
		return fmt.Errorf("prev is nil")
	}
	prev, next, locations, err := api.GetLocations(*cfg.prev)
	if err != nil {
		return err
	}
	cfg.next = next
	cfg.prev = prev
	for i := 0; i < len(locations); i++ {
		fmt.Println(locations[i])
	}
	return nil
}

func main() {
	commands := getCommands()

	cfg := &config{
		next: "https://pokeapi.co/api/v2/location/",
		prev: nil,
	}

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		reader.Scan()
		cmd := reader.Text()

		if cmd, ok := commands[cmd]; ok {
			cmd.callback(cfg)
		} else {
			fmt.Println("Command does not exist!")
		}
	}
}
