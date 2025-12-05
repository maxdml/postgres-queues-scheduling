package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"
)

// Export results to CSV file
func exportToCSV(tasks []Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"task_id", "duration_ms", "arrival_time", "dequeue_time",
		"completion_time", "wait_time_ms", "response_time_ms"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Calculate statistics
	var totalWaitTime, totalResponseTime time.Duration
	var minWaitTime, maxWaitTime, minResponseTime, maxResponseTime time.Duration
	waitTimes := make([]time.Duration, 0, len(tasks))
	responseTimes := make([]time.Duration, 0, len(tasks))
	firstTask := true

	// Write task data
	for _, task := range tasks {
		waitTime := task.DequeueTime.Sub(task.ArrivalTime)
		responseTime := task.CompletionTime.Sub(task.ArrivalTime)

		// Update statistics
		totalWaitTime += waitTime
		totalResponseTime += responseTime
		waitTimes = append(waitTimes, waitTime)
		responseTimes = append(responseTimes, responseTime)

		if firstTask {
			minWaitTime = waitTime
			maxWaitTime = waitTime
			minResponseTime = responseTime
			maxResponseTime = responseTime
			firstTask = false
		} else {
			if waitTime < minWaitTime {
				minWaitTime = waitTime
			}
			if waitTime > maxWaitTime {
				maxWaitTime = waitTime
			}
			if responseTime < minResponseTime {
				minResponseTime = responseTime
			}
			if responseTime > maxResponseTime {
				maxResponseTime = responseTime
			}
		}

		row := []string{
			fmt.Sprintf("%d", task.TaskID),
			fmt.Sprintf("%.0f", float64(task.Duration.Milliseconds())),
			task.ArrivalTime.Format(time.RFC3339Nano),
			task.DequeueTime.Format(time.RFC3339Nano),
			task.CompletionTime.Format(time.RFC3339Nano),
			fmt.Sprintf("%.3f", waitTime.Seconds()*1000),
			fmt.Sprintf("%.3f", responseTime.Seconds()*1000),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	fmt.Printf("\nResults exported to %s\n", filename)

	// Print summary statistics
	numTasks := len(tasks)
	if numTasks > 0 {
		// Sort times for percentile calculations
		sort.Slice(waitTimes, func(i, j int) bool { return waitTimes[i] < waitTimes[j] })
		sort.Slice(responseTimes, func(i, j int) bool { return responseTimes[i] < responseTimes[j] })

		// Calculate percentiles
		p90Index := int(float64(numTasks) * 0.90)
		if p90Index >= numTasks {
			p90Index = numTasks - 1
		}
		p99Index := int(float64(numTasks) * 0.99)
		if p99Index >= numTasks {
			p99Index = numTasks - 1
		}

		fmt.Printf("\nSummary Statistics:\n")
		/* print queueing delay if interested
		fmt.Printf("  Mean wait time: %.3f ms\n", float64(totalWaitTime.Milliseconds())/float64(numTasks))
		fmt.Printf("  Min wait time: %.3f ms\n", float64(minWaitTime.Milliseconds()))
		fmt.Printf("  Max wait time: %.3f ms\n", float64(maxWaitTime.Milliseconds()))
		fmt.Printf("  P90 wait time: %.3f ms\n", float64(waitTimes[p90Index].Milliseconds()))
		fmt.Printf("  P99 wait time: %.3f ms\n", float64(waitTimes[p99Index].Milliseconds()))
		*/
		fmt.Printf("  Mean response time: %.3f ms\n", float64(totalResponseTime.Milliseconds())/float64(numTasks))
		fmt.Printf("  Min response time: %.3f ms\n", float64(minResponseTime.Milliseconds()))
		fmt.Printf("  Max response time: %.3f ms\n", float64(maxResponseTime.Milliseconds()))
		fmt.Printf("  P90 response time: %.3f ms\n", float64(responseTimes[p90Index].Milliseconds()))
		fmt.Printf("  P99 response time: %.3f ms\n", float64(responseTimes[p99Index].Milliseconds()))
	}

	return nil
}
