package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

		if len(text) == 0 {
			continue
		}

		for k, c := range *commands {
			args := strings.Split(text, " ")
			if args[0] == k {
				if err := c.Callback(args...); err != nil {
					fmt.Println(err)
				}
				continue
			}
		}
	}
}
