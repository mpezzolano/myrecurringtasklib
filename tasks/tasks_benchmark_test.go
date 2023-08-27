package tasks

import (
	"testing"
	"time"
)

func BenchmarkAddTask(b *testing.B) {

	scheduler := New()

	taskFunc := func() error {
		return nil
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		task := &Task{
			Name:          "DummyTask",
			Interval:      1 * time.Second,
			RunOnce:       true,
			StartAfter:    time.Now(),
			TaskFunc:      taskFunc,
			ExecutionTime: 0,
			LastRun:       time.Time{},
			SuccessCount:  0,
			FailureCount:  0,
		}
		_, _ = scheduler.Add(task)
	}

	b.StopTimer()

	for name, _ := range scheduler.tasks {
		scheduler.Del(name)
	}
}
