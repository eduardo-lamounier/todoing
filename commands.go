package main

import (
	"errors"
	"fmt"
	"slices"
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
	TaskName string
}

func NewAddCommand(argument string) (*AddCommand, error) {
	c := AddCommand{argument}

	return &c, nil
}

type RemoveCommand struct {
	TaskNames []string
}

func NewRemoveCommand(arguments []string) *RemoveCommand {
	c := RemoveCommand{arguments}

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
		panic("Some flag must be specified to the root command.")
	}

	return nil
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

	selected, err := ScrollInput(len(priorities), func(currentPriority int) string {
		return fmt.Sprint("Priority: ", priorities[currentPriority].HumanText())
	}, true)

	return priorities[selected], err
}

func inputTaskState() (TaskState, error) {
	states := []TaskState{TaskPending, TaskInProgress, TaskCompleted}

	selected, err := ScrollInput(len(states), func(currentState int) string {
		return fmt.Sprint("State: ", states[currentState].HumanText())
	}, true)

	return states[selected], err
}

// TODO:
//   - dates should start from the current time
//   - there should have an option for not associating a deadline to a task
//     (as they're optional)
func inputTaskDeadline() (*time.Time, error) {
	now := time.Now()

	selectedYear, err := ScrollInput(
		now.Year()+10,
		func(currentYear int) string {
			return fmt.Sprintf("Deadline: %4v/",
				currentYear+1)
		},
		true,
	)
	if err != nil {
		return nil, err
	}

	selectedMonth, err := ScrollInput(12, func(currentMonth int) string {
		return fmt.Sprintf("Deadline: %4v/%2v/", selectedYear+1, currentMonth+1)
	}, false)
	if err != nil {
		return nil, err
	}

	// Assumes all months have 31 days for simplicity
	selectedDay, err := ScrollInput(DaysInMonth(time.Month(selectedMonth+1), selectedYear+1), func(currentDay int) string {
		return fmt.Sprintf("Deadline: %4v/%2v/%2v", selectedYear+1, selectedMonth+1, currentDay+1)
	}, false)
	if err != nil {
		return nil, err
	}

	selected := time.Date(selectedYear+1, time.Month(selectedMonth+1), selectedDay+1,
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
