package main

import (
	"os"

	"github.com/einarno/pokedexcli/commands"
)

func main() {
	commands.Start(os.Stdin, os.Stdout)
}
