package tasks

import (
	"fmt"
	"testing"
	"time"
)

var dummyTaskFunc = func() error {
	fmt.Println("Executing dummy task function.")
	return nil
}

var dummyErrFunc = func(err error) {
	fmt.Println("Executing dummy error function. Error: ", err)
}

func TestAdd(t *testing.T) {
	scheduler := New()
	task := &Task{
		Interval: time.Second,
		TaskFunc: &dummyTaskFunc,
		ErrFunc:  &dummyErrFunc,
	}

	id, err := scheduler.Add(task)

	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	if id == "" {
		t.Errorf("Expected ID to be non-empty")
	}
}

func TestDel(t *testing.T) {
	scheduler := New()
	task := &Task{
		Interval: time.Second,
		TaskFunc: &dummyTaskFunc,
		ErrFunc:  &dummyErrFunc,
	}

	id, err := scheduler.Add(task)

	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	scheduler.Del(id)
	_, err = scheduler.Lookup(id)

	if err == nil {
		t.Errorf("Task was not deleted.")
	}
}

func TestStop(t *testing.T) {
	scheduler := New()
	task1 := &Task{
		Interval: time.Second,
		TaskFunc: &dummyTaskFunc,
		ErrFunc:  &dummyErrFunc,
	}
	task2 := &Task{
		Interval: time.Second,
		TaskFunc: &dummyTaskFunc,
		ErrFunc:  &dummyErrFunc,
	}

	_, err := scheduler.Add(task1)
	if err != nil {
		return
	}
	_, err = scheduler.Add(task2)
	if err != nil {
		return
	}
	scheduler.Stop()

	if len(scheduler.tasks) != 0 {
		t.Errorf("Stop did not delete all tasks.")
	}
}
