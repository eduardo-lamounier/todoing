package main

import "time"

type Task struct {
	Name     string       `json:"name"`
	State    TaskState    `json:"category,omitempty"`
	Priority TaskPriority `json:"priority,omitempty"`
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

func (t *Task) WithState(state TaskState) *Task {
	t.State = state
	return t
}

func (t *Task) WithPriority(priority TaskPriority) *Task {
	t.Priority = priority
	return t
}

func (t *Task) WithDeadline(deadline time.Time) *Task {
	t.Deadline = &deadline
	return t
}

func (t *Task) WithNoDeadline() *Task {
	t.Deadline = nil
	return t
}

////////////////////////////////////////////////////////////////////////////////

func (s TaskState) String() string {
	return stateNames[s]
}

func (p TaskPriority) String() string {
	return priorityNames[p]
}
