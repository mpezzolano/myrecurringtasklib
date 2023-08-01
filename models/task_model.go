package models

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TaskError struct {
	ID    string
	Task  string
	Cause error
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("task %s (ID %s) failed: %v", e.Task, e.ID, e.Cause)
}

type Task struct {
	sync.Mutex
	id            string
	Name          string
	Interval      time.Duration
	RunOnce       bool
	StartAfter    time.Time
	TaskFunc      func() error
	OnFail        func(error)
	timer         *time.Timer
	ctx           context.Context
	cancel        context.CancelFunc
	ExecutionTime time.Duration
	LastRun       time.Time
	SuccessCount  int
	FailureCount  int
}

type Scheduler struct {
	sync.RWMutex
	tasks map[string]*Task
}
