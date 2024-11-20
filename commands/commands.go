package commands

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/einarno/pokedexcli/pokeapi"
	c "github.com/einarno/pokedexcli/pokecache"
	"github.com/einarno/pokedexcli/pokedata"
)

type config struct {
	offset  int
	cache   *c.Cache
	pokemap map[string]pokedata.Pokemon
}
type cliCommand struct {
	name        string
	description string
	callback    func(conf *config, args []string) error
}

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	io.WriteString(out, "Welcome to the Pokedex!\n")
	conf := config{offset: 0, cache: c.NewCache(10 * time.Minute), pokemap: make(map[string]pokedata.Pokemon)}
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}
		// Split input into command and arguments
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		cmdName := parts[0]
		args := parts[1:]

		cmdMap := getCliCommands(out)
		cmd, ok := cmdMap[cmdName]
		if ok {
			err := cmd.callback(&conf, args)
			if err != nil {
			}
		} else {
			io.WriteString(out, "Woops! looks like you used an unknown command\n")
		}
	}
}

func getCliCommands(out io.Writer) map[string]cliCommand {

	commandHelp := func(conf *config, _ []string) error {
		io.WriteString(out, "\nUsage:\n")
		for _, v := range getCliCommands(out) {
			commandString := fmt.Sprintf("\n%s: %s\n", v.name, v.description)
			io.WriteString(out, commandString)
		}
		io.WriteString(out, "\n\n")
		return nil
	}

	// I want to handle the exit with a gracefull return so this is just  a empty function
	commandExit := func(conf *config, _ []string) error {
		return nil
	}
	commandMap := func(conf *config, _ []string) error {
		places, err := pokeapi.GetLocationAreas(conf.cache, conf.offset)
		if err != nil {
			return err
		}
		for _, place := range places {
			commandString := fmt.Sprintf("%s\n", place)
			io.WriteString(out, commandString)
		}
		conf.offset += 20
		io.WriteString(out, fmt.Sprintf("%d\n", conf.offset))
		return nil
	}
	commandExplore := func(conf *config, args []string) error {
		areaID := args[0]
		places, err := pokeapi.ExploreArea(conf.cache, areaID)
		if err != nil {
			return err
		}
		io.WriteString(out, fmt.Sprintf("Exloring %s\n\nFound pokemon:\n", areaID))
		for _, place := range places {
			commandString := fmt.Sprintf("%s\n", place)
			io.WriteString(out, commandString)
		}
		return nil
	}
	commandCatch := func(conf *config, args []string) error {
		pokemonName := args[0]
		pokemon, err := pokeapi.GetPokemon(pokemonName)
		if err != nil {
			return err
		}
		io.WriteString(out, fmt.Sprintf("Throwing a Pokeball at %s...\n", pokemonName))
		isCatched := pokedata.IsCatched(*pokemon.BaseExperience)
		if isCatched {
			conf.pokemap[pokemonName] = pokemon
			io.WriteString(out, fmt.Sprintf("Caught %s\n", pokemonName))

		} else {
			io.WriteString(out, "Failed...\n")

		}
		return nil

	}
	commandPokedex := func(conf *config, args []string) error {
		for _, pokemon := range conf.pokemap {
			io.WriteString(out, pokemon.Name+"\n")
		}
		return nil
	}
	commandInspect := func(conf *config, args []string) error {
		pokemon, ok := conf.pokemap[args[0]]
		if !ok {
			io.WriteString(out, fmt.Sprintf("%s not caught\n", args[0]))
			return nil
		}
		printString := fmt.Sprintf("Name: %s\nHeight: %d\nWeight: %d\nStats:", pokemon.Name, pokemon.Height, pokemon.Weight)
		for _, stat := range pokemon.Stats {
			printString += fmt.Sprintf("\n  - %s: %d", stat.Stat.Name, stat.BaseStat)
		}
		printString += "\nTypes:"
		for _, pokemonType := range pokemon.Types {
			printString += fmt.Sprintf("\n  - %s", pokemonType.Type.Name)
		}
		io.WriteString(out, printString+"\n")
		return nil
	}
	commandMapb := func(conf *config, _ []string) error {
		if conf.offset < 20 {
			io.WriteString(out, "Error: No need to go back when you have not checked out the locations yet\n")
			return nil
		}
		places, err := pokeapi.GetLocationAreas(conf.cache, conf.offset)
		if err != nil {
			return err
		}
		for _, place := range places {
			commandString := fmt.Sprintf("%s\n", place)
			io.WriteString(out, commandString)
		}
		conf.offset -= 20
		io.WriteString(out, fmt.Sprintf("%d\n", conf.offset))
		return nil
	}
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
			description: "displays the names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the names of the last 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the pokemons in a certain area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch pokemon",
			callback:    commandCatch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See all pokemon",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "See pokemon details",
			callback:    commandInspect,
		},
	}
}
