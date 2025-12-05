package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

// Configuration constants
const (
	NUM_TASKS              = 100
	SHORT_TASK_DURATION    = 300 * time.Millisecond
	LONG_TASK_DURATION     = 2000 * time.Millisecond
	SHORT_TASK_PROBABILITY = 0.8
	TARGET_UTILIZATION     = 0.7
)

// Task represents a single task with timing information
type Task struct {
	TaskID         int
	Duration       time.Duration
	ArrivalTime    time.Time
	DequeueTime    time.Time
	CompletionTime time.Time
}

// TaskResult includes calculated metrics
type TaskResult struct {
	Task         Task
	WaitTime     time.Duration
	ResponseTime time.Duration
}

// Step to get current time (non-deterministic operation)
func getCurrentTime(ctx context.Context) (time.Time, error) {
	return time.Now(), nil
}

// Step to simulate work by sleeping
func simulateWork(_ context.Context, duration time.Duration) (string, error) {
	time.Sleep(duration)
	return "completed", nil
}

// Workflow to process a task
func processTask(ctx dbos.DBOSContext, task Task) (Task, error) {
	// Record dequeue time when workflow starts
	dequeueTime, err := dbos.RunAsStep(ctx, getCurrentTime)
	if err != nil {
		return task, err
	}
	task.DequeueTime = dequeueTime

	// Simulate work by sleeping for the task duration
	_, err = dbos.RunAsStep(ctx, func(stepCtx context.Context) (string, error) {
		return simulateWork(stepCtx, task.Duration)
	})
	if err != nil {
		return task, err
	}

	// Record completion time
	completionTime, err := dbos.RunAsStep(ctx, getCurrentTime)
	if err != nil {
		return task, err
	}
	task.CompletionTime = completionTime

	return task, nil
}

func main() {
	// Parse command-line flags
	algo := flag.String("algo", "fcfs", "Scheduling algorithm to use (fcfs, sjf)")
	flag.Parse()

	// Run the appropriate algorithm
	switch *algo {
	case "fcfs":
		FCFS()
	case "sjf":
		SJF()
	default:
		fmt.Printf("Unknown algorithm: %s\n", *algo)
		fmt.Println("Available algorithms: fcfs, sjf")
		os.Exit(1)
	}
}
