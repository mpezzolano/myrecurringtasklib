package tasks

import "fmt"

type TaskError struct {
	ID    string
	Task  string
	Cause error
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("task %s (ID %s) failed: %v", e.Task, e.ID, e.Cause)
}

type TaskNotFoundError struct {
	TaskID string
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with ID %s not found", e.TaskID)
}
