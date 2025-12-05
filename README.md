# Queue Scheduling Simulations

Queue scheduling is omnipresent in computer systems—from CPU schedulers to web servers to cloud job queues. This repository demonstrates situations where poor scheduling leads to performance degradation, for example, head of line blocking, and candidate solutions.

# Queueing theory background

I recommend obtaining a copy of [Mor Harchol-Balter's book](https://www.amazon.com/Performance-Modeling-Design-Computer-Systems/dp/1107027500/) and reading chapter I.2, "Queuing Theory Terminology"

# How the Simulator Works

The simulator models a queueing system where tasks arrive over time and are processed by a one or many worker. 
The simulator accepts parameters that define:
- The total number of tasks to simulate
- Task duration distributions (e.g., short vs. long task durations)
- The probability distribution for task types
- Target system utilization (how busy the server should be)

Right now tasks arrive at a uniform interval equal to their average processing time, weighed by the target utilization.

# The head of line blocking problem

Head of line blocking happens when a processor consumes tasks from a queue, and some tasks are much longer than others. This can result in a majority of (short) tasks being stuck behind a long task. This situation can be observed with a simple simulation: craft a workflow with a 80% ratio of short to long request, enqueue tasks as they arrive in a queue, and have the worker dequeue first, come, first serve (FCFS), the tasks.

# SJF vs FCFS

One potential solution is to prioritize short tasks over long tasks (a scheduling algorithm known as Shortest Job First). When giving higher priority to the short tasks, we ask the worker(s) to dequeue them first, even if some long tasks arrived in the queue earlier.

Here is an example with a workload made of 80% short tasks (100ms) and 20% long tasks (2 seconds).

Let's compare the latency distributions for all tasks, short then long tasks, on a simulation run. What we'll observe is that SJF massively improves performance for the short tasks, at the cost of a slight performance degradation for the long ones (specifically around the median).

![FCFS vs SJF Comparison](fcfs-sjf.png)

The risk with SJF is that we can totally starve the long tasks. In this case, we have a very gentle arrival distribution (requests arrive at the same interval in the system), but with more bursty arrivals, the likelihood of starving long jobs will augment.

# Try it yourself!

## Setup

This repo's queues are backed by Postgres. Get a Postgres instance and set your database connection:
```bash
export DBOS_SYSTEM_DATABASE_URL="postgresql://postgres:postgres@localhost:5432/queues"
```

Build:
```bash
go build -o main
```

## Running Experiments

Run FCFS (First Come First Served):
```bash
go run . -algo fcfs
```

Run SJF (Shortest Job First):
```bash
go run . -algo sjf
```

Each run generates a timestamped CSV file in the `results/` directory.

## Generating Plots

Compare the algorithms by plotting their results:
```bash
python plot_results.py [result.csv] [result.csv] ...
```

This generates `algorithm_comparison.png` showing average response time for each algorithm.