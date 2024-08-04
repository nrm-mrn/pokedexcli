package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/nrm-mrn/pokedexcli/internal/pokeAPI"
)

var configValues config

func startRepl() {
	configValues.baseExp = 300
	configValues.collection = make(map[string]pokeapi.Pokemon)

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
			default:
				configValues.param = &tokens[1]
				if command, ok := commandsBank[tokens[0]]; ok {
					err := command.callback(&configValues)
					if err != nil {
						fmt.Printf("%v\n", err)
					}
					configValues.param = nil
				} else {
					unknownCommand()
				}
			}
			continue
		}
		if len(tokens) > 2 {
			unknownCommand()
			continue
		}
		if len(tokens) == 0 {
			continue
		}
		if command, ok := commandsBank[tokens[0]]; ok {
			err := command.callback(&configValues)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			if command.name == "exit" {
				os.Exit(0)
			}
		} else {
			unknownCommand()
		}
	}
}
