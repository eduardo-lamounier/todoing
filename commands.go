package main

import "fmt"

const helpMessage = `usage: todoing [-v | --version] [-h | --help] <command> [args]

NOTE: The below commands are not yet implemented.

list: Lists all the added tasks.

add: Adds a new task.

remove: Removes the specified task.`

const version = "1.0.0"

type Command interface {
	Execute()
}

type RootCommand struct {
	flag string
}

func (c RootCommand) Execute() {
	switch c.flag {
	case "--help", "-h":
		fmt.Println(helpMessage)
	case "--version", "-v":
		fmt.Println("todoing", version)
	default:
		fmt.Printf("Unknown flag '%s'.\n", c.flag)
	}
}
