package main

import "fmt"

const helpMessage = `usage: todoing [-v | --version] [-h | --help] <command> [args]

NOTE: The below commands are not yet implemented.

list: Lists all the added tasks.

add: Adds a new task.

remove: Removes the specified task.`

const version = "1.0.0"

type Command interface {
	Execute() error
}

type RootCommand struct {
	flag string
}

func (c RootCommand) Execute() error {
	switch c.flag {
	case "--help", "-h":
		fmt.Println(helpMessage)
	case "--version", "-v":
		fmt.Println("todoing", version)
	default:
		return fmt.Errorf("unknown flag '%s'", c.flag)
	}

	return nil
}
