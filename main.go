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
		throw(err)
		return
	}

	err = command.Execute()
	if err != nil {
		throw(err)
	}
}

func throw(err error) {
	fmt.Printf("ERROR: %s.\n", err)
	os.Exit(1)
}

func parse(args []string) (Command, error) {
	if len(args) == 0 {
		return nil,
			errors.New("no command, flags or options passed to the program")
	}

	switch args[0] {
	case "list":
		return NewListCommand(), nil
	case "add":
		if len(args[1:]) < 1 {
			return nil, fmt.Errorf("command 'add' expects an argument")
		}

		if len(args[1:]) > 1 {
			return nil, fmt.Errorf("too many arguments passed to command 'add'")
		}

		return NewAddCommand(args[1])
	case "remove":
		return NewRemoveCommand(), nil
	default:
		if len(args[0]) < 2 || !strings.HasPrefix(args[0], "-") {
			return nil, fmt.Errorf("unknown command '%s'", args[0])
		}

		return NewRootCommand(args[0])
	}
}
