package main

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
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

type AddCommand struct {
	CommandName string
}

func NewAddCommand(argument string) (*AddCommand, error) {
	c := AddCommand{argument}

	return &c, nil
}

type RemoveCommand struct {
	CommandName string
}

func NewRemoveCommand(argument string) *RemoveCommand {
	c := RemoveCommand{argument}

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

	colsName := []string{
		"Name:",
		"Priority:",
		"Deadline:",
		"State:",
	}

	var rows [][]string
	for _, task := range userData.Tasks {
		row := []string{
			task.Name,
			formatPriority(task.Priority),
			formatDeadline(task.Deadline),
			formatState(task.State),
		}

		rows = append(rows, row)
	}

	err = PrintTable("[ Tasks ]", "You have no tasks.", rows, colsName)
	if err != nil {
		log.Fatal("error when listing your tasks:", err)
	}
	return nil
}

func (c AddCommand) Execute() error {
	task := NewTask(c.CommandName)

	userData, err := LoadUserData()
	if err != nil {
		return err
	}

	userData.Tasks = append(userData.Tasks, *task)
	return SaveUserData(userData)
}

func (c RemoveCommand) Execute() error {
	userData, err := LoadUserData()
	if err != nil {
		return err
	}

	taskIdx := -1
	for i, task := range userData.Tasks {
		if task.Name == c.CommandName {
			taskIdx = i
			break
		}
	}

	if taskIdx == -1 {
		return fmt.Errorf("task with name '%s' does not exist", c.CommandName)
	}

	userData.Tasks = slices.Delete(userData.Tasks, taskIdx, taskIdx+1)
	return SaveUserData(userData)
}
