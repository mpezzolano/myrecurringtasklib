package main

import (
	"fmt"
	"github.com/mpezzolano/myrecurringtasklib/tasks"

	"time"
)

func main() {
	scheduler := tasks.New()

	// Def. task 1
	taskFunc1 := func() error {
		for i := 1; i <= 10; i++ {
			fmt.Println(i)
		}
		return nil
	}
	task1 := &tasks.Task{
		Interval:   time.Second * 5,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   &taskFunc1,
	}

	// Def. task 2
	taskFunc2 := func() error {
		fmt.Println("Hello World")
		return nil
	}
	task2 := &tasks.Task{
		Interval:   time.Second * 5,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   &taskFunc2,
	}

	// Add tasks to scheduler
	_, _ = scheduler.Add(task1)
	_, _ = scheduler.Add(task2)

	time.Sleep(time.Second * 10)
	fmt.Println("finish tasks")
}
