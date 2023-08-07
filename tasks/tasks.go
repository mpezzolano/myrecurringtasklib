package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/xid"
	"go.uber.org/zap"
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

var logger, _ = zap.NewProduction()

func New() *Scheduler {
	logger.Info("Creating new scheduler")
	return &Scheduler{
		tasks: make(map[string]*Task),
	}
}

func (schd *Scheduler) Add(t *Task) (string, error) {
	logger.Info("Adding a new task")
	id := xid.New()
	err := schd.AddWithID(id.String(), t)
	if err != nil {
		logger.Error("Failed to add task", zap.Error(err))
	}
	return id.String(), err
}

func (schd *Scheduler) AddWithID(id string, t *Task) error {
	if t.TaskFunc == nil {
		logger.Error("Task function cannot be nil")
		return fmt.Errorf("task function cannot be nil")
	}
	if t.Interval <= time.Duration(0) {
		logger.Error("Task interval must be defined")
		return fmt.Errorf("task interval must be defined")
	}

	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.id = id
	schd.Lock()
	defer schd.Unlock()

	if _, ok := schd.tasks[id]; ok {
		logger.Error("ID already used")
		return fmt.Errorf("ID already used")
	}

	schd.tasks[t.id] = t
	schd.scheduleTask(t)
	return nil
}

func (schd *Scheduler) Del(name string) {
	logger.Info("Deleting task", zap.String("task", name))
	t, err := schd.Lookup(name)
	if err != nil {
		logger.Error("Failed to delete task", zap.Error(err))
		return
	}

	t.Lock()
	defer t.Unlock()

	if t.timer != nil {
		t.timer.Stop()
	}
	t.cancel()
	schd.Lock()
	defer schd.Unlock()
	delete(schd.tasks, name)
}

func (schd *Scheduler) Lookup(name string) (*Task, error) {
	schd.RLock()
	defer schd.RUnlock()

	t, ok := schd.tasks[name]
	if ok {
		return t, nil
	}
	logger.Error("Could not find task within the task list", zap.String("task", name))
	return nil, fmt.Errorf("could not find task within the task list")
}

func (schd *Scheduler) Stop() {
	logger.Info("Stopping scheduler")
	schd.Lock()
	defer schd.Unlock()

	for n, t := range schd.tasks {
		t.cancel()
		delete(schd.tasks, n)
	}
}

func (schd *Scheduler) scheduleTask(t *Task) {
	_ = time.AfterFunc(time.Until(t.StartAfter), func() {
		if err := t.ctx.Err(); err == nil {
			t.Lock()
			t.timer = time.AfterFunc(t.Interval, func() { schd.execTask(t) })
			t.Unlock()
		}
	})
}

func (schd *Scheduler) execTask(t *Task) {
	go func() {
		startTime := time.Now()
		err := (t.TaskFunc)()
		executionTime := time.Since(startTime)

		t.Lock()
		t.ExecutionTime = executionTime
		t.LastRun = startTime
		if err != nil && t.OnFail != nil {
			t.FailureCount++
			go t.OnFail(&TaskError{
				ID:    t.id,
				Task:  t.Name,
				Cause: err,
			})
			logger.Error("Task failed", zap.String("task", t.Name), zap.Error(err))
		} else {
			t.SuccessCount++
		}
		t.Unlock()

		if t.RunOnce {
			schd.Del(t.id)
		} else {
			t.Lock()
			t.timer.Reset(t.Interval)
			t.Unlock()
		}
	}()
}
