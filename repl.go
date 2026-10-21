package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <location-area>")
	}

	locationName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	fmt.Printf("Exploring %s...\n", locationName)

	area, err := getLocationAreaPokemon(url)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandMap(cfg *config, args ...string) error {
	url := locationAreaURL

	if cfg.next != nil {
		url = *cfg.next
	}

	areas, err := getLocationAreas(url)
	if err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	cfg.next = areas.Next
	cfg.previous = areas.Previous

	return nil

}

func commandMapBack(cfg *config, args ...string) error {
	if cfg.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	areas, err := getLocationAreas(*cfg.previous)
	if err != nil {
		return err
	}

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}

	cfg.next = areas.Next
	cfg.previous = areas.Previous

	return nil

}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	next     *string
	previous *string
	commands map[string]cliCommand
	pokedex  map[string]Pokemon
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range cfg.commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	return nil

}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			return
		}

		words := cleanInput(scanner.Text())

		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		cmd, ok := cfg.commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.callback(cfg, words[1:]...)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon>")
	}

	name := strings.ToLower(args[0])

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := getPokemon(name)
	if err != nil {
		return err
	}

	// Higher base experience means a harder catch.
	catchChance := 100 - pokemon.BaseExperience/10

	if catchChance < 10 {
		catchChance = 10
	}
	if catchChance > 90 {
		catchChance = 90
	}

	if rand.Intn(100) >= catchChance {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)

	if cfg.pokedex == nil {
		cfg.pokedex = make(map[string]Pokemon)
	}

	cfg.pokedex[pokemon.Name] = pokemon

	return nil

}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}

	name := strings.ToLower(args[0])

	pokemon, ok := cfg.pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}

	return nil

}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Println("Your Pokedex:")

	for _, pokemon := range cfg.pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil

}
