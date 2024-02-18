package app

import (
	"bufio"
	"fmt"
	"os"

	"github.com/unclesamrocks/pokedexcli/internal/cli"
)

func Init() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := cli.New()

	fmt.Println("Welcome to the Pokedex!")
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		text := scanner.Text()

		if command, ok := (*commands)[text]; ok {
			if err := command.Callback(); err != nil {
				fmt.Println(err)
			}

		} else {
			continue
		}
	}
}
