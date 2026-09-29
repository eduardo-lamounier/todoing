package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Name     string       `json:"name"`
	State    TaskState    `json:"category"`
	Priority TaskPriority `json:"priority"`
	Deadline *time.Time   `json:"deadline,omitempty"`
}

func NewTask(name string) *Task {
	task := Task{
		Name:     name,
		State:    TaskPending,
		Priority: TaskMediumPriority,
		Deadline: nil,
	}

	return &task
}

type (
	TaskState    int
	TaskPriority int
)

const (
	TaskPending TaskState = iota
	TaskInProgress
	TaskCompleted
)

const (
	TaskLowPriority TaskPriority = iota
	TaskMediumPriority
	TaskHighPriority
	TaskUrgentPriority
)

var stateNames = map[TaskState]string{
	TaskPending:    "pending",
	TaskInProgress: "in_progress",
	TaskCompleted:  "completed",
}

var priorityNames = map[TaskPriority]string{
	TaskLowPriority:    "low",
	TaskMediumPriority: "medium",
	TaskHighPriority:   "high",
	TaskUrgentPriority: "urgent",
}

////////////////////////////////////////////////////////////////////////////////

func (s TaskState) String() string {
	return stateNames[s]
}

func (p TaskPriority) String() string {
	return priorityNames[p]
}

func (s TaskState) HumanText() string {
	return Capitalize(strings.ReplaceAll(s.String(), "_", " "))
}

func (p TaskPriority) HumanText() string {
	return Capitalize(p.String())
}

func (s TaskState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (p TaskPriority) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (s *TaskState) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	for state, name := range stateNames {
		if name == str {
			*s = state
			return nil
		}
	}

	return fmt.Errorf("invalid byte sequence to convert to TaskState")
}

func (p *TaskPriority) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	for priority, name := range priorityNames {
		if name == str {
			*p = priority
			return nil
		}
	}

	return fmt.Errorf("invalid byte sequence to convert to TaskPriority")
}
