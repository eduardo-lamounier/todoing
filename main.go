package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	command, err := parse(args)
	if err != nil {
		fmt.Printf("ERROR: %s.", err)
		return
	}

	command.Execute()
}

func parse(args []string) (Command, error) {
	if len(args) == 0 {
		return nil,
			errors.New("no command, flags or options passed to the program")
	}

	switch args[0] {
	default:
		if len(args[0]) < 2 || !strings.HasPrefix(args[0], "-") {
			return nil, fmt.Errorf("unknown command '%s'", args[0])
		}

		return RootCommand{args[0]}, nil
	}
}
