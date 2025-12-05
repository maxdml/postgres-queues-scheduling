package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/dbos-inc/dbos-transact-golang/dbos"
)

// FCFS implements the First-Come-First-Served scheduling algorithm
func FCFS() {
	avgTaskDuration := time.Duration(float64(SHORT_TASK_DURATION)*SHORT_TASK_PROBABILITY +
		float64(LONG_TASK_DURATION)*(1-SHORT_TASK_PROBABILITY))
	interArrivalTime := time.Duration(float64(avgTaskDuration) / TARGET_UTILIZATION)

	fmt.Println("============================================================")
	fmt.Println("FCFS: First-Come-First-Served Queue Scheduling Demo")
	fmt.Println("============================================================")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Number of tasks: %d\n", NUM_TASKS)
	fmt.Printf("  Short task duration: %v\n", SHORT_TASK_DURATION)
	fmt.Printf("  Long task duration: %v\n", LONG_TASK_DURATION)
	fmt.Printf("  Short task probability: %.0f%%\n", SHORT_TASK_PROBABILITY*100)
	fmt.Printf("  Average task duration: %v\n", avgTaskDuration)
	fmt.Printf("  Target utilization: %.0f%%\n", TARGET_UTILIZATION*100)
	fmt.Printf("  Average inter-arrival time: %v\n", interArrivalTime)
	fmt.Printf("  Queue: Single FIFO queue with single worker\n")
	fmt.Println("============================================================")

	// Initialize DBOS context with PostgreSQL
	dbosContext, err := dbos.NewDBOSContext(context.Background(), dbos.Config{
		AppName:     "fifo-queue-demo",
		DatabaseURL: os.Getenv("DBOS_SYSTEM_DATABASE_URL"),
	})
	if err != nil {
		panic(fmt.Sprintf("Initializing DBOS failed: %v", err))
	}

	// Create a single FIFO queue with worker concurrency of 1 (single worker)
	fifoQueue := dbos.NewWorkflowQueue(dbosContext, "fifo_queue", dbos.WithWorkerConcurrency(1), dbos.WithQueueBasePollingInterval(100*time.Millisecond), dbos.WithQueueMaxPollingInterval(100*time.Millisecond))

	// Register the workflow
	dbos.RegisterWorkflow(dbosContext, processTask)

	// Launch DBOS
	err = dbos.Launch(dbosContext)
	if err != nil {
		panic(fmt.Sprintf("Launching DBOS failed: %v", err))
	}
	defer dbos.Shutdown(dbosContext, 5*time.Second)

	// Enqueue tasks one at a time, respecting arrival times
	fmt.Printf("\nEnqueueing tasks to FIFO queue with respect to arrival times...\n")
	startTime := time.Now()
	handles := make([]dbos.WorkflowHandle[Task], NUM_TASKS)
	completedTasks := make([]Task, NUM_TASKS)
	shortCount := 0
	longCount := 0

	for i := range NUM_TASKS {
		// Pick task duration based on probability
		var duration time.Duration
		if rand.Float64() < SHORT_TASK_PROBABILITY {
			duration = SHORT_TASK_DURATION
			shortCount++
		} else {
			duration = LONG_TASK_DURATION
			longCount++
		}

		// Calculate arrival time for this task
		expectedArrivalTime := startTime.Add(time.Duration(i) * interArrivalTime)

		// Sleep until the task is due
		now := time.Now()
		if expectedArrivalTime.After(now) {
			time.Sleep(expectedArrivalTime.Sub(now))
		}

		// Create task with current time as arrival time
		task := Task{
			TaskID:      i,
			Duration:    duration,
			ArrivalTime: time.Now(),
		}

		// Enqueue the task
		handle, err := dbos.RunWorkflow(dbosContext, processTask, task, dbos.WithQueue(fifoQueue.Name))
		if err != nil {
			panic(fmt.Sprintf("Failed to enqueue task %d: %v", i, err))
		}
		handles[i] = handle

		if (i+1)%10 == 0 {
			fmt.Printf("  Enqueued %d/%d tasks...\n", i+1, NUM_TASKS)
		}
	}

	fmt.Printf("\nAll %d tasks enqueued (%d short, %d long). Processing...\n", NUM_TASKS, shortCount, longCount)

	// Wait for all tasks to complete and collect results
	for i, handle := range handles {
		result, err := handle.GetResult()
		if err != nil {
			panic(fmt.Sprintf("Task %d failed: %v", i, err))
		}
		completedTasks[i] = result
		if (i+1)%10 == 0 {
			fmt.Printf("  Completed %d/%d tasks...\n", i+1, NUM_TASKS)
		}
	}

	fmt.Printf("\nAll %d tasks completed!\n", len(completedTasks))

	// Create results directory if it doesn't exist
	resultsDir := "results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create results directory: %v", err))
	}

	// Generate unique filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(resultsDir, fmt.Sprintf("fifo_results_%s.csv", timestamp))

	// Export results to CSV
	fmt.Printf("\nExporting results...\n")
	if err := exportToCSV(completedTasks, filename); err != nil {
		panic(fmt.Sprintf("Failed to export CSV: %v", err))
	}

	fmt.Println("\n============================================================")
	fmt.Println("Demo completed successfully!")
	fmt.Println("============================================================")
}

