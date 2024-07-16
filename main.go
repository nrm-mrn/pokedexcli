package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Printf("Here are all of the available commands:\n")
	commandsBank := getCommands()
	for key := range commandsBank {
		fmt.Printf("%v\n", key)
	}
	fmt.Printf("To get more info about a command use help <command>\n")
	return nil
}
func commandExit() error {
	fmt.Printf("Goodbuy!")
	return nil
}

func getCommands() map[string]cliCommand {
	helpCom := cliCommand{
		name:        "help",
		description: "displays a help message\n",
		callback:    commandHelp,
	}
	exitCom := cliCommand{
		name:        "exit",
		description: "exits the application\n",
		callback:    commandExit,
	}
	return map[string]cliCommand{
		"help": helpCom,
		"exit": exitCom,
	}
}

func unknownCommand() {
	fmt.Printf("Unknown command, consult help for a list of available commands\n")
}
func main() {
	fmt.Printf("Welcome to pokedexCLI tool!\n")
	commandsBank := getCommands()
	for {
		fmt.Printf("Gimme your command, sir!\n")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userInput := scanner.Text()
		tokens := strings.Split(userInput, " ")
		if len(tokens) == 2 {
			switch tokens[0] {
			case "help":
				if command, ok := commandsBank[tokens[1]]; ok {
					fmt.Printf(command.description)
				} else {
					unknownCommand()
				}
			}
			continue
		}
		if len(tokens) > 2 {
			unknownCommand()
		}
		if command, ok := commandsBank[tokens[0]]; ok {
			command.callback()
			if command.name == "exit" {
				break
			}
		} else {
			unknownCommand()
		}
	}

}
