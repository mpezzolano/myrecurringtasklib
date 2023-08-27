// task_errors.go

package tasks

import "fmt"

type TaskNotFoundError struct {
	TaskID string
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("Task with ID %s not found", e.TaskID)
}

type TaskExecutionError struct {
	TaskID  string
	Message string
}

func (e *TaskExecutionError) Error() string {
	return fmt.Sprintf("Error executing task with ID %s: %s", e.TaskID, e.Message)
}
