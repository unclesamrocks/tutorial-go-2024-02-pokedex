package cli

import (
	"fmt"
	"os"
)

type cliCommand struct {
	Name        string
	Description string
	Callback    func() error
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

	return &commands
}

func createCallbackHelp(commands *map[string]cliCommand) func() error {
	fn := func() error {
		fmt.Printf("\nUsage:\n\n")
		for name, command := range *commands {
			fmt.Printf("%s: %s\n", name, command.Description)
		}
		fmt.Println()

		return nil
	}

	return fn
}

func callbackExit() error {
	os.Exit(0)
	return nil
}
