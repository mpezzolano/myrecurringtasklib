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

	// failing
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
		TaskFunc:   &printNumbers,
		OnFail:     onFail,
	})

	schd.Add(&tasks.Task{
		Name:       "Print Hello World",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   &printHelloWorld,
		OnFail:     onFail,
	})

	schd.Add(&tasks.Task{
		Name:       "Failing Task",
		Interval:   time.Second * 1,
		RunOnce:    true,
		StartAfter: time.Now(),
		TaskFunc:   &failingTask,
		OnFail:     onFail,
	})

	time.Sleep(time.Second * 5)
	fmt.Println("Finish")
}
