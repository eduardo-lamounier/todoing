package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

const helpMessage = `usage: todoing [-v | --version] [-h | --help] <command> [args]

Supported commands:

  list      Lists all the current tasks.

  add       Adds a new task with the specified name.

  remove    Removes the task with the specified name.

You can use the [ -h | --help ] flag in any command to see its detailed reference.`

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

type ConfigCommand struct {
	ShowHelp bool
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

type ListCommand struct {
	ShowHelp bool
}

func NewListCommand(flags []string) (*ListCommand, error) {
	c := ListCommand{
		ShowHelp: slices.Contains(flags, "-h") || slices.Contains(flags, "--help"),
	}

	var errs []error
	for _, flag := range flags {
		if len(flag) < 2 || !strings.HasPrefix(flag, "-") {
			continue
		}

		errs = append(
			errs,
			fmt.Errorf("unknown flag '%s' for command 'list'", flag),
		)
	}

	return &c, errors.Join(errs...)
}

type AddCommand struct {
	ShowHelp bool

	TaskName string
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

type RemoveCommand struct {
	ShowHelp bool

	TaskNames []string
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
	userData, err := LoadUserData()
	if err != nil {
		return err
	}

	formatDeadline := func(deadline *time.Time) string {
		if deadline == nil {
			return "N/A"
		}
		return deadline.Format("01/02/2006")
	}

	colsName := []string{
		"Name:",
		"Priority:",
		"Deadline:",
		"State:",
	}

	var rows [][]string
	rows = slices.Grow(rows, len(userData.Tasks))

	for _, task := range userData.Tasks {
		row := []string{
			task.Name,
			task.Priority.HumanText(),
			formatDeadline(task.Deadline),
			task.State.HumanText(),
		}

		rows = append(rows, row)
	}

	PrintTable("[ Tasks ]", "You have no tasks.", rows, colsName)
	return nil
}

func inputTaskPriority() (TaskPriority, error) {
	priorities := []TaskPriority{
		TaskMediumPriority,
		TaskHighPriority,
		TaskUrgentPriority,
		TaskLowPriority,
	}

	selected, err := WrapPrompt(-1, 0, len(priorities), func(currentPriority int) string {
		return fmt.Sprint("Priority: ", priorities[currentPriority].HumanText())
	}, true)

	return priorities[selected], err
}

func inputTaskState() (TaskState, error) {
	states := []TaskState{TaskPending, TaskInProgress, TaskCompleted}

	selected, err := WrapPrompt(-1, 0, len(states), func(currentState int) string {
		return fmt.Sprint("State: ", states[currentState].HumanText())
	}, true)

	return states[selected], err
}

func inputTaskDeadline() (*time.Time, error) {
	shouldInputDeadline, err := YesNoPrompt(
		"Would you like to set a deadline for this task?", false, true,
	)

	if !shouldInputDeadline || err != nil {
		return nil, err
	}

	now := time.Now()

	selectedYear, err := WrapPrompt(
		-1, now.Year(), now.Year()+10+1,
		func(currentYear int) string {
			return fmt.Sprintf("Deadline [YYYY/MM/DD]: %4v/", currentYear)
		},
		false,
	)
	if err != nil {
		return nil, err
	}

	selectedMonth, err := WrapPrompt(-1, int(now.Month()), 12+1,
		func(currentMonth int) string {
			return fmt.Sprintf("Deadline [YYYY/MM/DD]: %4v/%2v/", selectedYear, currentMonth)
		}, false)
	if err != nil {
		return nil, err
	}

	maxDay := DaysInMonth(time.Month(selectedMonth), selectedYear)
	selectedDay, err := WrapPrompt(min(now.Day(), maxDay), 1, maxDay+1,
		func(currentDay int) string {
			return fmt.Sprintf("Deadline [YYYY/MM/DD]: %4v/%2v/%2v",
				selectedYear, selectedMonth, currentDay)
		}, true)
	if err != nil {
		return nil, err
	}

	selected := time.Date(selectedYear, time.Month(selectedMonth), selectedDay,
		0, 0, 0, 0, now.Location())
	return &selected, nil
}

func (c AddCommand) Execute() error {
	userData, err := LoadUserData()
	if err != nil {
		return err
	}

	for _, task := range userData.Tasks {
		if task.Name == c.TaskName {
			return fmt.Errorf("task \"%s\" already exists", task.Name)
		}
	}

	const (
		priorityInputText = "Priority: "
		deadlineInputText = "Deadline: "
		stateInputText    = "State: "
	)

	fmt.Printf("[ %s ]\n", c.TaskName)

	task := NewTask(c.TaskName)

	task.Priority, err = inputTaskPriority()
	if err != nil {
		return err
	}
	task.State, err = inputTaskState()
	if err != nil {
		return err
	}
	task.Deadline, err = inputTaskDeadline()
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

	// A set. Stores the tasks yet to be removed (a key is the task's name)
	toRemove := make(map[string]struct{})

	for _, taskName := range c.TaskNames {
		toRemove[taskName] = struct{}{}
	}

	for i, task := range userData.Tasks {
		if _, shouldRemove := toRemove[task.Name]; shouldRemove {
			userData.Tasks = slices.Delete(userData.Tasks, i, i+1)
			delete(toRemove, task.Name)
		}
	}

	var errs []error
	for taskName := range toRemove {
		// Any task remaining in `toRemove` weren't found in the user's data:
		errs = append(errs, fmt.Errorf("task with name \"%s\" does not exist", taskName))
	}

	if err := SaveUserData(userData); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
