package tasks

import (
	"context"
	"sync"
	"time"
)

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
