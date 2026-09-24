package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No command, flags or options passed to the program.")
		return
	}

	var command Command
	switch args[0] {
	default:
		if len(args[0]) < 2 || !strings.HasPrefix(args[0], "-") {
			fmt.Printf("Unknown command '%s'.\n", args[0])
			return
		}

		command = RootCommand{args[0]}
	}

	command.Execute()
}
