package main

import (
	"bufio"
	"fmt"
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
	}
}

func commandHelp() error {
	commands := getCommands()
	commandOrder := []string{"help", "exit"}
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

func commandExit() error {
	os.Exit(0)
	return nil
}

func main() {
	commands := getCommands()

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		reader.Scan()
		cmd := reader.Text()

		if cmd, ok := commands[cmd]; ok {
			cmd.callback()
		} else {
			fmt.Println("Command does not exist!")
		}
	}
}
