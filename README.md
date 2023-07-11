# Task Scheduler in Go

The tasks package in this repository is a versatile task scheduler written in Go. It allows you to schedule tasks at regular intervals, or run tasks just once at a specific time. Moreover, it gives you control over handling task errors.

## How it Works

A task in this context is a function with the signature func() error. This function represents the task that needs to be performed. The task is wrapped in a struct, Task, which includes several configuration options such as the interval to run the task, the time after which to start the task, and an error handling function in case the task encounters an error during execution.

The main component of the tasks package is the Scheduler. The Scheduler is responsible for managing and executing all tasks. You can add tasks to a Scheduler, and they will be executed in the background at their configured intervals.
## Using the Task Scheduler

### Creating a New Scheduler

To create a new Scheduler, simply call the `New()` function from the `tasks` package:

```
scheduler := tasks.New()
```

This will return a new Scheduler instance.

### Defining a Task

A task is a function that matches the TaskFunc type (i.e., func() error). For example:

```
taskFunc := func() error {
	fmt.Println("Hello, world!")
	return nil
}
```

### Scheduling a Task

To schedule a task, you must first create a Task instance, which includes the task function and various 
scheduling options:

```
task := &tasks.Task{
    Interval:   time.Second * 5,
    RunOnce:    true,
    StartAfter: time.Now(),
    TaskFunc:   &taskFunc,
}
```

Then, you can add the task to the scheduler:

```
_, err := scheduler.Add(task)
if err != nil {
	// handle error
}
```

Running the Scheduler

Once you have added all the tasks you want to schedule, you can run the scheduler:

```
scheduler.Start()
```

That's it! The scheduler will now start running your tasks in the background at their specified intervals.

### Handling Errors

The tasks package comes with enhanced error handling. The OnFail function in each task allows for specific error handling to be defined per task. The package also defines custom error types for more granular error handling and debugging.

### Contributing

This is a detailed README with explanations on how the task scheduler works, and examples on how to use it. You may want to customize it to match the specifics of your project.
