package main

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

func ListTasks() error {
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

func AddTaskInteractively(taskName string) error {
	userData, err := LoadUserData()
	if err != nil {
		return err
	}

	for _, task := range userData.Tasks {
		if task.Name == taskName {
			return fmt.Errorf("task \"%s\" already exists", task.Name)
		}
	}

	const (
		priorityInputText = "Priority: "
		deadlineInputText = "Deadline: "
		stateInputText    = "State: "
	)

	fmt.Printf("[ %s ]\n", taskName)

	task := NewTask(taskName)

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

func RemoveTasks(taskNames []string, removeAll bool) error {
	userData, err := LoadUserData()
	if err != nil {
		return err
	}

	// A set. Stores the tasks yet to be removed (a key is the task's name)
	toRemove := make(map[string]struct{})

	for _, taskName := range taskNames {
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

	if removeAll {
		userData.Tasks = []Task{}
		fmt.Println("Removed all your tasks.")
	}

	if err := SaveUserData(userData); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func UpdateStateOfTasks(taskNames []string, state TaskState) error {
	userData, err := LoadUserData()
	if err != nil {
		return nil
	}

	// A set. Stores the tasks yet to have their states updated
	// (a key is the task's name)
	toUpdate := make(map[string]struct{})

	for _, taskName := range taskNames {
		toUpdate[taskName] = struct{}{}
	}

	for i, task := range userData.Tasks {
		if _, shouldUpdate := toUpdate[task.Name]; shouldUpdate {
			userData.Tasks[i].State = state
			delete(toUpdate, task.Name)
		}
	}

	var errs []error
	for taskName := range toUpdate {
		errs = append(errs, fmt.Errorf("task with name \"%s\" does not exist", taskName))
	}

	if err := SaveUserData(userData); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
