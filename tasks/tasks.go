package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/xid"
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
	id         string
	Name       string
	Interval   time.Duration
	RunOnce    bool
	StartAfter time.Time
	TaskFunc   func() error
	OnFail     func(error)
	timer      *time.Timer
	ctx        context.Context
	cancel     context.CancelFunc
}

type Scheduler struct {
	sync.RWMutex
	tasks map[string]*Task
}

func New() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
	}
}

func (schd *Scheduler) Add(t *Task) (string, error) {
	id := xid.New()
	return id.String(), schd.AddWithID(id.String(), t)
}

func (schd *Scheduler) AddWithID(id string, t *Task) error {
	if t.TaskFunc == nil {
		return fmt.Errorf("task function cannot be nil")
	}
	if t.Interval <= time.Duration(0) {
		return fmt.Errorf("task interval must be defined")
	}

	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.id = id
	schd.Lock()
	defer schd.Unlock()

	if _, ok := schd.tasks[id]; ok {
		return fmt.Errorf("ID already used")
	}

	schd.tasks[t.id] = t
	schd.scheduleTask(t)
	return nil
}

func (schd *Scheduler) Del(name string) {
	t, err := schd.Lookup(name)
	if err != nil {
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
	return nil, fmt.Errorf("could not find task within the task list")
}

func (schd *Scheduler) Stop() {
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
		if err := (t.TaskFunc)(); err != nil && t.OnFail != nil {
			go t.OnFail(&TaskError{
				ID:    t.id,
				Task:  t.Name,
				Cause: err,
			})
		}
		if t.RunOnce {
			schd.Del(t.id)
		} else {
			t.Lock()
			t.timer.Reset(t.Interval)
			t.Unlock()
		}
	}()
}
