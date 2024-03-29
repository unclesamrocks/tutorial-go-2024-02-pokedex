package cli

import (
	"fmt"
	"os"

	"github.com/unclesamrocks/pokedexcli/internal/api"
)

type cliCommand struct {
	Name        string
	Description string
	Callback    func(...string) error
}

func New() *map[string]cliCommand {
	commands := map[string]cliCommand{}

	commands["help"] = cliCommand{
		Name:        "help",
		Description: "Displays a help message",
		Callback:    createCallbackHelp(&commands),
	}

	commands["exit"] = cliCommand{
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    callbackExit,
	}

	// api commands
	api := api.New()

	commands["map"] = cliCommand{
		Name:        "map",
		Description: "Displays the names of 20 next location areas in the Pokemon world",
		Callback:    api.FetchNext,
	}

	commands["mapb"] = cliCommand{
		Name:        "mapb",
		Description: "Displays the names of 20 prev location areas in the Pokemon world",
		Callback:    api.FetchPrev,
	}

	commands["explore"] = cliCommand{
		Name:        "explore",
		Description: "Explore specific area of the world",
		Callback:    api.Explore,
	}

	commands["catch"] = cliCommand{
		Name:        "catch",
		Description: "Catch pokemon!",
		Callback:    api.Catch,
	}

	commands["inspect"] = cliCommand{
		Name:        "inspect",
		Description: "Inspect pokemon!",
		Callback:    api.Inspect,
	}

	commands["pokedex"] = cliCommand{
		Name:        "pokedex",
		Description: "Display all caught pokemon currently in Pokedex",
		Callback:    api.Pokedex,
	}

	return &commands
}

func createCallbackHelp(commands *map[string]cliCommand) func(...string) error {
	fn := func(args ...string) error {
		fmt.Printf("\nUsage:\n\n")
		for name, command := range *commands {
			fmt.Printf("%s: %s\n", name, command.Description)
		}
		fmt.Println()

		return nil
	}

	return fn
}

func callbackExit(args ...string) error {
	os.Exit(0)
	return nil
}
