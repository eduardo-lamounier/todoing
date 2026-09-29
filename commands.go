package main

import (
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

	selectedMonth, err := ScrollInput(12, func(currentMonth int) string {
		return fmt.Sprintf("Deadline: %2v/", currentMonth+1)
	}, false)
	if err != nil {
		return nil, err
	}

	// Assumes all months have 31 days for simplicity
	selectedDay, err := ScrollInput(31, func(currentDay int) string {
		return fmt.Sprintf("Deadline: %2v/%2v", selectedMonth+1, currentDay+1)
	}, false)
	if err != nil {
		return nil, err
	}

	selectedYear, err := ScrollInput(
		now.Year()+10,
		func(currentYear int) string {
			return fmt.Sprintf("Deadline: %2v/%2v/%4v",
				selectedMonth+1, selectedDay+1, currentYear+1)
		},
		true,
	)
	if err != nil {
		return nil, err
	}

	selected := time.Date(selectedYear+1, time.Month(selectedMonth+1), selectedDay+1,
		0, 0, 0, 0, now.Location())
	return &selected, nil
}

// Execute TODO: Do not allow the user to add a task with a non-unique name
func (c AddCommand) Execute() error {
	const (
		priorityInputText = "Priority: "
		deadlineInputText = "Deadline: "
		stateInputText    = "State: "
	)

	fmt.Printf("[ %s ]\n", c.CommandName)

	var err error

	task := NewTask(c.CommandName)

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
