package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// WorkloadConfig holds the workload configuration parameters
type WorkloadConfig struct {
	NumTasks             int     `yaml:"num_tasks"`
	ShortTaskDurationMs  int     `yaml:"short_task_duration_ms"`
	LongTaskDurationMs   int     `yaml:"long_task_duration_ms"`
	ShortTaskProbability float64 `yaml:"short_task_probability"`
	TargetUtilization    float64 `yaml:"target_utilization"`
}

// Config holds all application configuration
type Config struct {
	Workload WorkloadConfig `yaml:"workload"`
}

// Global configuration instance
var AppConfig Config

// LoadConfig loads configuration from config.yaml file
// If the file doesn't exist or has missing values, it uses defaults
func LoadConfig() error {
	// Set defaults
	AppConfig = Config{
		Workload: WorkloadConfig{
			NumTasks:             100,
			ShortTaskDurationMs:  300,
			LongTaskDurationMs:   2000,
			ShortTaskProbability: 0.8,
			TargetUtilization:    0.7,
		},
	}

	// Try to read config file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		// If file doesn't exist, use defaults
		if os.IsNotExist(err) {
			fmt.Println("No config.yaml found, using default configuration")
			return nil
		}
		return fmt.Errorf("failed to read config.yaml: %w", err)
	}

	// Parse YAML
	var fileConfig Config
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	// Merge file config with defaults (only override non-zero values)
	if fileConfig.Workload.NumTasks > 0 {
		AppConfig.Workload.NumTasks = fileConfig.Workload.NumTasks
	}
	if fileConfig.Workload.ShortTaskDurationMs > 0 {
		AppConfig.Workload.ShortTaskDurationMs = fileConfig.Workload.ShortTaskDurationMs
	}
	if fileConfig.Workload.LongTaskDurationMs > 0 {
		AppConfig.Workload.LongTaskDurationMs = fileConfig.Workload.LongTaskDurationMs
	}
	if fileConfig.Workload.ShortTaskProbability > 0 {
		AppConfig.Workload.ShortTaskProbability = fileConfig.Workload.ShortTaskProbability
	}
	if fileConfig.Workload.TargetUtilization > 0 {
		AppConfig.Workload.TargetUtilization = fileConfig.Workload.TargetUtilization
	}

	fmt.Println("Configuration loaded from config.yaml")
	return nil
}

// Helper methods to get durations as time.Duration
func (c *WorkloadConfig) ShortTaskDuration() time.Duration {
	return time.Duration(c.ShortTaskDurationMs) * time.Millisecond
}

func (c *WorkloadConfig) LongTaskDuration() time.Duration {
	return time.Duration(c.LongTaskDurationMs) * time.Millisecond
}

