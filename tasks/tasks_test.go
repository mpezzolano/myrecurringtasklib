package tasks

import (
	"errors"
	"testing"
	"time"
)

func mockTaskFunc() error {
	return nil
}

func mockFailingTaskFunc() error {
	return errors.New("test error")
}

func TestAddTask(t *testing.T) {
	scheduler := New()

	taskID, err := scheduler.Add(&Task{
		Name:     "TestTask",
		Interval: time.Second * 1,
		TaskFunc: mockTaskFunc,
	})

	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	task, _ := scheduler.Lookup(taskID)
	if task.Name != "TestTask" {
		t.Errorf("Task was not correctly added to the scheduler")
	}
}

func TestDeleteTask(t *testing.T) {
	scheduler := New()

	taskID, _ := scheduler.Add(&Task{
		Name:     "TestTask",
		Interval: time.Second * 1,
		TaskFunc: mockTaskFunc,
	})

	scheduler.Del(taskID)
	task, err := scheduler.Lookup(taskID)

	if err == nil {
		t.Errorf("Task was not correctly deleted: %v", task)
	}
}

func TestTaskExecution(t *testing.T) {
	scheduler := New()

	taskID, _ := scheduler.Add(&Task{
		Name:     "TestTask",
		Interval: time.Second * 1,
		TaskFunc: mockTaskFunc,
	})

	time.Sleep(time.Second * 2)

	task, _ := scheduler.Lookup(taskID)
	if task.SuccessCount != 1 {
		t.Errorf("Task did not execute successfully")
	}
}

func TestFailingTask(t *testing.T) {
	scheduler := New()

	taskID, _ := scheduler.Add(&Task{
		Name:     "FailingTask",
		Interval: time.Second * 1,
		TaskFunc: mockFailingTaskFunc,
	})

	time.Sleep(time.Second * 2)

	task, _ := scheduler.Lookup(taskID)
	if task.FailureCount != 1 {
		t.Errorf("Task failure was not registered")
	}
}

func TestTaskError(t *testing.T) {
	err := &TaskError{
		ID:    "1",
		Task:  "FailingTask",
		Cause: errors.New("test error"),
	}

	expectedError := "task FailingTask (ID 1) failed: test error"

	if err.Error() != expectedError {
		t.Errorf("Got %v, want %v", err.Error(), expectedError)
	}
}

func TestStopScheduler(t *testing.T) {
	scheduler := New()

	scheduler.Add(&Task{
		Name:     "TestTask",
		Interval: time.Second * 1,
		TaskFunc: mockTaskFunc,
	})

	scheduler.Stop()

	if len(scheduler.tasks) != 0 {
		t.Errorf("Scheduler did not stop all tasks")
	}
}
