package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const helpMessage = `usage: todoing [-v | --version] [-h | --help] <command> [args]

Supported commands:

  list               Lists all the current tasks.

  add                Adds a new task with the specified name.

  remove             Removes the tasks with the specified names.

  reset              Marks the tasks with the specified names as pending.

  begin              Marks the tasks with the specified names as in progress.

  complete           Marks the tasks with the specified names as completed.

Flags:

  [-h | --help]      Shows this help message.

  [-v | --version]   Shows the version of the app that is installed.


You can use the [ -h | --help ] flag in any command to see its detailed reference.`

const listHelpMessage = `usage: todoing list [-h | --help]

Lists all the tasks and their information: priority, state, deadline etc.

Flags:
	[-h | --help]   Shows this help message.`

const addHelpMessage = `usage: todoing add [-h | --help] [task name]

Adds a task with name of the specified argument.

After running this command, you'll be asked to enter the other informations
about the task iteratively (some are optional).

Flags:
	[-h | --help]   Shows this help message.`

const removeHelpMessage = `usage: todoing remove [-h | --help] [task names]

Removes the task with name of the specified argument.

Flags:
  [-h | --help]   Shows this help message.`

const resetHelpMessage = `usage: todoing reset [-h | --help] [task names]

Marks the tasks with specified names as pending.

Flags:
	[-h | --help]   Shows this help message.`

const beginHelpMessage = `usage: todoing begin [-h | --help] [task names]

Marks the tasks with specified names as in progress.

Flags:
	[-h | --help]   Shows this help message.`

const completeHelpMessage = `usage: todoing complete [-h | --help] [task names]

Marks the tasks with specified names as completed.

Flags:
	[-h | --help]   Shows this help message.`

const version = "1.0.0"

type Command interface {
	Execute() error
}

type RootCommand struct {
	ShowHelp    bool
	ShowVersion bool
}

type ConfigCommand struct {
	ShowHelp bool
}

type ListCommand struct {
	ShowHelp bool
}

type AddCommand struct {
	ShowHelp bool

	TaskName string
}

type RemoveCommand struct {
	ShowHelp bool

	TaskNames []string
}

type ResetCommand struct {
	ShowHelp bool

	TaskNames []string
}

type BeginCommand struct {
	ShowHelp bool

	TaskNames []string
}

type CompleteCommand struct {
	ShowHelp bool

	TaskNames []string
}

////////////////////////////////////////////////////////////////////////////////

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

func NewConfigCommand(flag string) (*ConfigCommand, error) {
	c := ConfigCommand{
		ShowHelp: flag == "-h" || flag == "--help",
	}

	if !c.ShowHelp {
		return nil, fmt.Errorf("unknown flag '%s' for command 'config'", flag)
	}

	return &c, nil
}

func NewListCommand(flags []string) (*ListCommand, error) {
	c := ListCommand{
		ShowHelp: slices.Contains(flags, "-h") || slices.Contains(flags, "--help"),
	}

	validFlags := map[string]struct{}{
		"-h": {}, "--help": {},
	}

	var errs []error
	for _, flag := range flags {
		if len(flag) < 2 || !strings.HasPrefix(flag, "-") {
			continue
		}

		if _, isValid := validFlags[flag]; isValid {
			continue
		}

		errs = append(
			errs,
			fmt.Errorf("unknown flag '%s' for command 'list'", flag),
		)
	}

	return &c, errors.Join(errs...)
}

func NewAddCommand(param string) (*AddCommand, error) {
	var c AddCommand
	if len(param) < 2 || !strings.HasPrefix(param, "-") {
		c = AddCommand{false, param}
		return &c, nil
	}

	c = AddCommand{
		ShowHelp: param == "-h" || param == "--help",
	}

	if !c.ShowHelp {
		return nil, fmt.Errorf("unknown flag '%s' for command 'add'", param)
	}

	return &c, nil
}

func NewRemoveCommand(params []string) (*RemoveCommand, error) {
	c := RemoveCommand{
		ShowHelp: slices.Contains(params, "-h") || slices.Contains(params, "--help"),
	}

	if c.ShowHelp {
		return &c, nil
	}

	var errs []error

	for _, param := range params {
		if len(param) >= 2 && strings.HasPrefix(param, "-") {
			errs = append(errs,
				fmt.Errorf("unknown flag '%s' for command 'remove'", param))
		} else {
			c.TaskNames = append(c.TaskNames, param)
		}
	}

	return &c, errors.Join(errs...)
}

func NewResetCommand(params []string) (*ResetCommand, error) {
	c := ResetCommand{
		ShowHelp: slices.Contains(params, "-h") || slices.Contains(params, "--help"),
	}

	validFlags := map[string]struct{}{
		"-h": {}, "--help": {},
	}

	var errs []error
	for _, param := range params {
		if len(param) < 2 || !strings.HasPrefix(param, "-") {
			c.TaskNames = append(c.TaskNames, param)
			continue
		}

		if _, isValid := validFlags[param]; isValid {
			continue
		}

		errs = append(
			errs,
			fmt.Errorf("unknown flag '%s' for command 'reset'", param),
		)
	}

	return &c, errors.Join(errs...)
}

func NewBeginCommand(params []string) (*BeginCommand, error) {
	c := BeginCommand{
		ShowHelp: slices.Contains(params, "-h") || slices.Contains(params, "--help"),
	}

	validFlags := map[string]struct{}{
		"-h": {}, "--help": {},
	}

	var errs []error
	for _, param := range params {
		if len(param) < 2 || !strings.HasPrefix(param, "-") {
			c.TaskNames = append(c.TaskNames, param)
			continue
		}

		if _, isValid := validFlags[param]; isValid {
			continue
		}

		errs = append(
			errs,
			fmt.Errorf("unknown flag '%s' for command 'begin'", param),
		)
	}

	return &c, errors.Join(errs...)
}

func NewCompleteCommand(params []string) (*CompleteCommand, error) {
	c := CompleteCommand{
		ShowHelp: slices.Contains(params, "-h") || slices.Contains(params, "--help"),
	}

	validFlags := map[string]struct{}{
		"-h": {}, "--help": {},
	}

	var errs []error
	for _, param := range params {
		if len(param) < 2 || !strings.HasPrefix(param, "-") {
			c.TaskNames = append(c.TaskNames, param)
			continue
		}

		if _, isValid := validFlags[param]; isValid {
			continue
		}

		errs = append(
			errs,
			fmt.Errorf("unknown flag '%s' for command 'complete'", param),
		)
	}

	return &c, errors.Join(errs...)
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
		panic("Some flag must have been specified to the root command.")
	}

	return nil
}

func (c ConfigCommand) Execute() error {
	return fmt.Errorf("command 'config' not yet implemented")
}

func (c ListCommand) Execute() error {
	if c.ShowHelp {
		fmt.Println(listHelpMessage)
		return nil
	}

	return ListTasks()
}

func (c AddCommand) Execute() error {
	if c.ShowHelp {
		fmt.Println(addHelpMessage)
		return nil
	}

	return AddTaskInteractively(c.TaskName)
}

func (c RemoveCommand) Execute() error {
	if c.ShowHelp {
		fmt.Println(removeHelpMessage)
		return nil
	}

	return RemoveTasks(c.TaskNames)
}

func (c ResetCommand) Execute() error {
	if c.ShowHelp {
		fmt.Println(resetHelpMessage)
		return nil
	}

	return fmt.Errorf("command 'reset' not yet implemented")
}

func (c BeginCommand) Execute() error {
	if c.ShowHelp {
		fmt.Println(beginHelpMessage)
		return nil
	}

	return fmt.Errorf("command 'begin' not yet implemented")
}

func (c CompleteCommand) Execute() error {
	if c.ShowHelp {
		fmt.Println(completeHelpMessage)
		return nil
	}

	return fmt.Errorf("command 'complete' not yet implemented")
}
