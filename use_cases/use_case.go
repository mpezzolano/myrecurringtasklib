package main

import (
	"errors"
	"fmt"
	"github.com/mpezzolano/myrecurringtasklib/tasks"
	"time"
)

func main() {
	printNumbers := func() error {
		for i := 1; i <= 10; i++ {
			fmt.Println(i)
		}
		return nil
	}

	printHelloWorld := func() error {
		fmt.Println("Hello World")
		return nil
	}

	// failing task
	failingTask := func() error {
		return errors.New("intentional failure")
	}

	onFail := func(err error) {
		fmt.Printf("Error: %v\n", err)
	}

	schd := tasks.New()

	schd.Add(&tasks.Task{
		Name:       "Print Numbers",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   printNumbers,
		OnFail:     onFail,
	})

	schd.Add(&tasks.Task{
		Name:       "Print Hello World",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   printHelloWorld,
		OnFail:     onFail,
	})

	schd.Add(&tasks.Task{
		Name:       "Failing Task",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   failingTask,
		OnFail:     onFail,
	})

	// Try to lookup a task that doesn't exist
	_, err := schd.Lookup("NonExistingTask")
	if err != nil {
		// We'll check if the error is of type TaskNotFoundError
		if _, ok := err.(*tasks.TaskNotFoundError); ok {
			fmt.Println("Error Custom: task doesn't exist")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}

	time.Sleep(time.Second * 5)
	fmt.Println("Finish")
}
