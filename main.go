package main

func main() {
	cfg := &config{
		commands: map[string]cliCommand{
			"pokedex": {
				name:        "pokedex",
				description: "List all caught Pokemon",
				callback:    commandPokedex,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspect a caught Pokemon",
				callback:    commandInspect,
			},
			"catch": {
				name:        "catch",
				description: "Catch a Pokemon",
				callback:    commandCatch,
			},
			"explore": {
				name:        "explore",
				description: "Explore a location area",
				callback:    commandExplore,
			},
			"map": {
				name:        "map",
				description: "Display the next 20 location areas",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Display the previous 20 location areas",
				callback:    commandMapBack,
			},
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
		},
	}

	startRepl(cfg)

}
