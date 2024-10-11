package commands

import (
	"bufio"
	"fmt"
	"io"

	"github.com/einarno/pokedexcli/pokeapi"
)

type config struct {
	offset int
}
type cliCommand struct {
	name        string
	description string
	callback    func(conf *config) error
}

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	io.WriteString(out, "Welcome to the Pokedex!\n")
	conf := config{offset: 0}
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
		cmdMap := getCliCommands(out)
		cmd, ok := cmdMap[line]
		if ok {
			err := cmd.callback(&conf)
			if err != nil {

			}
		} else {
			io.WriteString(out, "Woops! looks like you used an unknown command\n")
		}
	}
}

func getCliCommands(out io.Writer) map[string]cliCommand {

	commandHelp := func(conf *config) error {
		io.WriteString(out, "\nUsage:\n")
		for _, v := range getCliCommands(out) {
			commandString := fmt.Sprintf("\n%s: %s\n", v.name, v.description)
			io.WriteString(out, commandString)
		}
		io.WriteString(out, "\n\n")
		return nil
	}

	// I want to handle the exit with a gracefull return so this is just  a empty function
	commandExit := func(conf *config) error {
		return nil
	}
	commandMap := func(conf *config) error {
		places, err := pokeapi.GetLocationAreas(conf.offset)
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
	commandMapb := func(conf *config) error {
		if conf.offset < 20 {
			io.WriteString(out, "Error: No need to go back when you have not checked out the locations yet\n")
			return nil
		}
		places, err := pokeapi.GetLocationAreas(conf.offset)
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
	}
}
