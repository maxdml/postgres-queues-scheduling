package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Load configuration
	if err := LoadConfig(); err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

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
