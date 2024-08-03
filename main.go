package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/nrm-mrn/pokedexcli/internal/pokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}
type config struct {
	nextArea *string
	prevArea *string
	param    *string
}

var configValues config

func commandHelp(c *config) error {
	fmt.Printf("Here are all of the available commands:\n")
	commandsBank := getCommands()
	for key := range commandsBank {
		fmt.Printf("%v\n", key)
	}
	fmt.Printf("To get more info about a command use help <command>\n")
	return nil
}
func commandExit(c *config) error {
	fmt.Printf("Goodbuy!")
	return nil
}

func commandMap(c *config) error {
	var locationsStruct pokeapi.LocationsResponse
	var err error
	locationsStruct, err = pokeapi.GetLocations(c.nextArea)
	if err != nil {
		return err
	}

	c.nextArea = locationsStruct.Next
	c.prevArea = locationsStruct.Previous
	for _, area := range locationsStruct.Results {
		fmt.Printf("%v\n", area.Name)
	}
	return nil

}

func commandMapb(c *config) error {
	if c.prevArea == nil {
		return errors.New("Already on the first page")
	}
	locationsStruct, err := pokeapi.GetLocations(c.prevArea)
	if err != nil {
		return err
	}
	c.nextArea = locationsStruct.Next
	c.prevArea = locationsStruct.Previous
	for _, area := range locationsStruct.Results {
		fmt.Printf("%v\n", area.Name)
	}
	return nil
}

func commandExplore(c *config) error {
	if c.param == nil {
		return errors.New("Area name is nil, use map for list of all names")
	}
	exploredLoca, err := pokeapi.ExploreLocation(c.param)
	if err != nil {
		return err
	}
	for _, encounter := range exploredLoca.PokemonEncounters {
		fmt.Printf("%v\n", encounter.Pokemon.Name)
	}
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
	mapCom := cliCommand{
		name:        "map",
		description: "displays 20 location areas in the pokemon world. Each subsequent call to map will display the next 20 locations and so on",
		callback:    commandMap,
	}
	mapbCom := cliCommand{
		name:        "mapb",
		description: "displays 20 previous location areas if previous page exists\n",
		callback:    commandMapb,
	}
	exploreCom := cliCommand{
		name:        "explore <location_name>",
		description: "displays pokemon names that one can encounter in specified location\n",
		callback:    commandExplore,
	}
	return map[string]cliCommand{
		"help":    helpCom,
		"exit":    exitCom,
		"map":     mapCom,
		"mapb":    mapbCom,
		"explore": exploreCom,
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
			default:
				configValues.param = &tokens[1]
				if command, ok := commandsBank[tokens[0]]; ok {
					err := command.callback(&configValues)
					if err != nil {
						fmt.Printf("%v\n", err)
					}
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
