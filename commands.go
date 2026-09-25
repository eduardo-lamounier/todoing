package main

import (
	"errors"
	"fmt"
)

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
	ShowHelp    bool
	ShowVersion bool
}

func NewRootCommand(flag string) (*RootCommand, error) {
	c := RootCommand{
		ShowHelp:    flag == "--help" || flag == "-h",
		ShowVersion: flag == "--version" || flag == "-v",
	}

	if !c.ShowHelp && !c.ShowVersion {
		return nil, fmt.Errorf("unknown flag '%s'", flag)
	}

	return &c, nil
}

type ListCommand struct{}

func NewListCommand() *ListCommand {
	c := ListCommand{}

	return &c
}

type AddCommand struct{}

func NewAddCommand() *AddCommand {
	c := AddCommand{}

	return &c
}

type RemoveCommand struct{}

func NewRemoveCommand() *RemoveCommand {
	c := RemoveCommand{}

	return &c
}

////////////////////////////////////////////////////////////////////////////////

func (c RootCommand) Execute() error {
	switch {
	case c.ShowHelp:
		fmt.Println(helpMessage)
	case c.ShowVersion:
		fmt.Println("todoing", version)
	default:
		// Shouldn't reach here
		return errors.New("some flag is required when no command is passed")
	}

	return nil
}

func (c ListCommand) Execute() error {
	return errors.New("command 'list' not yet implemented")
}

func (c AddCommand) Execute() error {
	return errors.New("command 'add' not yet implemented")
}

func (c RemoveCommand) Execute() error {
	return errors.New("command 'remove' not yet implemented")
}
