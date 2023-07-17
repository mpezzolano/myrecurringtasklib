package tasks

import (
	"errors"
	"testing"
	"time"
)

func TestScheduler_Add(t *testing.T) {
	s := New()

	taskFunc := func() error {
		return nil
	}

	task := &Task{
		Name:       "Test Task",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   taskFunc,
	}

	taskID, err := s.Add(task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	if taskID == "" {
		t.Errorf("Expected task ID to be non-empty")
	}
}

func TestScheduler_AddWithID(t *testing.T) {
	s := New()

	taskFunc := func() error {
		return nil
	}

	task := &Task{
		Name:       "Test Task",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   taskFunc,
	}

	err := s.AddWithID("testID", task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	_, err = s.Lookup("testID")
	if err != nil {
		t.Errorf("Failed to lookup task: %v", err)
	}
}

func TestTaskError_Error(t *testing.T) {
	taskErr := &TaskError{
		ID:    "testID",
		Task:  "Test Task",
		Cause: errors.New("test error"),
	}

	expectedErrorMsg := "task Test Task (ID testID) failed: test error"
	if taskErr.Error() != expectedErrorMsg {
		t.Errorf("Expected error message to be '%s', got '%s'", expectedErrorMsg, taskErr.Error())
	}
}

func TestTask_Fail(t *testing.T) {
	s := New()

	taskFunc := func() error {
		return errors.New("test error")
	}

	task := &Task{
		Name:       "Test Task",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   taskFunc,
		OnFail: func(err error) {
			expectedErrorMsg := "task Test Task (ID testID) failed: test error"
			if err.Error() != expectedErrorMsg {
				t.Errorf("Expected error message to be '%s', got '%s'", expectedErrorMsg, err.Error())
			}
		},
	}

	err := s.AddWithID("testID", task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}
}
