package main

import (
	"errors"
	"fmt"
	"math/rand"

	pokeapi "github.com/nrm-mrn/pokedexcli/internal/pokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}
type config struct {
	nextArea   *string
	prevArea   *string
	param      *string
	baseExp    int
	collection map[string]pokeapi.Pokemon
}

var client pokeapi.Client = pokeapi.NewClient()

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
	locationsStruct, err = pokeapi.GetLocations(&client, c.nextArea)
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
	locationsStruct, err := pokeapi.GetLocations(&client, c.prevArea)
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
		return errors.New("Area name is nil, use map for list of all names\n")
	}
	exploredLoca, err := pokeapi.ExploreLocation(&client, c.param)
	if err != nil {
		return err
	}
	for _, encounter := range exploredLoca.PokemonEncounters {
		fmt.Printf("%v\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config) error {
	if c.param == nil {
		return errors.New("Pokemon name is nil, use explore for possible names\n")
	}
	_, ok := c.collection[*c.param]
	if ok {
		fmt.Printf("%s is already caught and in collection\n", *c.param)
		return nil
	}
	pokemonInfo, err := pokeapi.GetPokemon(&client, *c.param)
	if err != nil {
		return err
	}
	catchProbability := 1 - float64(pokemonInfo.BaseExperience)/float64(c.baseExp)
	hit := rand.Float64()
	if hit < catchProbability {
		fmt.Printf(
			"Pokemon %s caught with hit %.2f and probability %.2f\nYou may now inspect it with inspect command\n",
			*c.param,
			hit,
			catchProbability,
		)
		c.collection[*c.param] = pokemonInfo
	} else {
		fmt.Printf("Pokemon %s NOT caught with hit %.2f and probability %.2f\n", *c.param, hit, catchProbability)
	}
	return nil
}

func commandInspect(c *config) error {
	if c.param == nil {
		if len(c.collection) == 0 {
			fmt.Printf("Your collection is empty\n")
		} else {
			fmt.Printf("Your current collection:\n")
			for _, p := range c.collection {
				fmt.Printf("- %s\n", p.Name)
			}
		}
		return errors.New("Pokemon name is nil, use one from your collection above\n")
	}
	p, ok := c.collection[*c.param]
	if !ok {
		return errors.New(fmt.Sprintf("%s is not in your collection, catch it first", *c.param))
	}

	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", p.Name, p.Height, p.Weight)
	fmt.Printf("Stats:\n")
	for _, s := range p.Stats {
		fmt.Printf(" -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range p.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(c *config) error {
	if len(c.collection) == 0 {
		fmt.Printf("Your collection is empty\n")
	} else {
		fmt.Printf("Your current collection:\n")
		for _, p := range c.collection {
			fmt.Printf("- %s\n", p.Name)
		}
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
	catchCom := cliCommand{
		name:        "catch <pokemon>",
		description: "attempts to catch a <pokemon>.\n",
		callback:    commandCatch,
	}
	inspectCom := cliCommand{
		name:        "inspect <pokemon>",
		description: "shows description of pokemon stats. Works only on caught pokemons.",
		callback:    commandInspect,
	}
	pokedexCom := cliCommand{
		name:        "pokedex",
		description: "lists all of your caught pokemons",
		callback:    commandPokedex,
	}
	return map[string]cliCommand{
		"help":    helpCom,
		"exit":    exitCom,
		"map":     mapCom,
		"mapb":    mapbCom,
		"explore": exploreCom,
		"catch":   catchCom,
		"inspect": inspectCom,
		"pokedex": pokedexCom,
	}
}

func unknownCommand() {
	fmt.Printf("Unknown command, consult help for a list of available commands\n")
}
