package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
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
	userData, err := LoadUserData()
	if err != nil {
		return err
	}

	const (
		minWidth = 0
		tabWidth = 4
		padding  = 4
		padChar  = ' '
		flags    = 0
	)

	tabWriter := tabwriter.NewWriter(
		os.Stdout, minWidth, tabWidth, padding, padChar, flags,
	)

	formatPriority := func(priority TaskPriority) string {
		return Capitalize(priority.String())
	}

	formatDeadline := func(deadline *time.Time) string {
		if deadline == nil {
			return "N/A"
		}
		return deadline.Format("01/02/2006")
	}

	formatState := func(state TaskState) string {
		return Capitalize(strings.ReplaceAll(state.String(), "_", " "))
	}

	widthCols := [4]int{
		len("Name:"), len("Priority:"), len("Deadline:"), len("State:"),
	}

	for _, task := range userData.Tasks {
		widthCols[0] = max(widthCols[0], len(task.Name))
		widthCols[1] = max(widthCols[1], len(formatPriority(task.Priority)))
		widthCols[2] = max(widthCols[2], len(formatDeadline(task.Deadline)))
		widthCols[3] = max(widthCols[3], len(formatState(task.State)))
	}

	totalWidth := (len(widthCols) - 1) * padding
	for _, widthCol := range widthCols {
		totalWidth += widthCol
	}

	const title = "[ Tasks ]"
	fmt.Printf("%*s\n\n", totalWidth/2+len(title)/2, title)
	fmt.Fprintf(tabWriter, "Name:\tPriority:\tDeadline:\tState:\n")
	for _, task := range userData.Tasks {
		fmt.Fprintf(
			tabWriter, "%v\t%v\t%v\t%v\n",
			task.Name,
			formatPriority(task.Priority),
			formatDeadline(task.Deadline),
			formatState(task.State),
		)
	}
	tabWriter.Flush()

	return nil
}

func (c AddCommand) Execute() error {
	return errors.New("command 'add' not yet implemented")
}

func (c RemoveCommand) Execute() error {
	return errors.New("command 'remove' not yet implemented")
}
