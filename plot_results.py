#!/usr/bin/env python3
"""
Plot response time comparison for DBOS queue scheduling results.

Compares median, p90, and p99 response times across different scheduling algorithms.
Response time = completion_time - arrival_time
"""

import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import sys
from pathlib import Path
import re


def extract_algorithm_name(filename):
    """Extract algorithm name from filename pattern 'algo_results_timestamp.csv'"""
    stem = Path(filename).stem
    # Match pattern: algo_results_timestamp
    match = re.match(r'^([^_]+)_results_', stem)
    if match:
        return match.group(1).upper()
    return stem


def calculate_percentiles(series):
    """Calculate median, p90, and p99 for a pandas Series."""
    return {
        'median': series.median(),
        'p90': series.quantile(0.90),
        'p99': series.quantile(0.99)
    }


def identify_task_types(df):
    """Identify short and long tasks based on duration_ms values."""
    # Get unique duration values and identify short vs long
    unique_durations = sorted(df['duration_ms'].unique())
    
    if len(unique_durations) >= 2:
        # Assume shortest is short, longest is long
        short_duration = min(unique_durations)
        long_duration = max(unique_durations)
    elif len(unique_durations) == 1:
        # Only one duration type, treat as short
        short_duration = unique_durations[0]
        long_duration = None
    else:
        return None, None
    
    short_tasks = df[df['duration_ms'] == short_duration]
    long_tasks = df[df['duration_ms'] == long_duration] if long_duration is not None else pd.DataFrame()
    
    return short_tasks, long_tasks


def plot_response_time_comparison(csv_files):
    """Generate bar charts comparing median, p90, and p99 response times."""
    
    # Collect data from all CSV files
    results = []
    
    for csv_file in csv_files:
        if not Path(csv_file).exists():
            print(f"Warning: File '{csv_file}' not found! Skipping...")
            continue
            
        # Read the CSV file
        df = pd.read_csv(csv_file)
        
        # Extract algorithm name
        algo_name = extract_algorithm_name(csv_file)
        
        # Calculate statistics for all tasks
        all_stats = calculate_percentiles(df['response_time_ms'])
        
        # Identify short and long tasks
        short_tasks, long_tasks = identify_task_types(df)
        
        # Calculate statistics for short tasks
        short_stats = None
        if short_tasks is not None and len(short_tasks) > 0:
            short_stats = calculate_percentiles(short_tasks['response_time_ms'])
        
        # Calculate statistics for long tasks
        long_stats = None
        if long_tasks is not None and len(long_tasks) > 0:
            long_stats = calculate_percentiles(long_tasks['response_time_ms'])
        
        results.append({
            'algorithm': algo_name,
            'all': all_stats,
            'short': short_stats,
            'long': long_stats,
            'file': csv_file
        })
        
        print(f"Loaded {csv_file}: {algo_name}")
        print(f"  All tasks - Median: {all_stats['median']:.2f} ms, P90: {all_stats['p90']:.2f} ms, P99: {all_stats['p99']:.2f} ms")
        if short_stats:
            print(f"  Short tasks - Median: {short_stats['median']:.2f} ms, P90: {short_stats['p90']:.2f} ms, P99: {short_stats['p99']:.2f} ms")
        if long_stats:
            print(f"  Long tasks - Median: {long_stats['median']:.2f} ms, P90: {long_stats['p90']:.2f} ms, P99: {long_stats['p99']:.2f} ms")
    
    if not results:
        print("Error: No valid CSV files found!")
        sys.exit(1)
    
    # Sort results by algorithm name for consistent display
    results = sorted(results, key=lambda x: x['algorithm'])
    
    # Extract algorithm names
    algorithms = [r['algorithm'] for r in results]
    num_algorithms = len(algorithms)
    
    # Create figure with 2 subplots in a column (Short and Long tasks only)
    fig, axes = plt.subplots(2, 1, figsize=(14, 12))
    
    # Define metrics and their positions
    metrics = ['Median', 'P90', 'P99']
    x = np.arange(num_algorithms)
    width = 0.25  # Width of bars
    
    # Colors for each metric
    colors = {
        'Median': '#2E86AB',  # Blue
        'P90': '#A23B72',     # Purple
        'P99': '#F18F01'      # Orange
    }
    
    # Plot 1: All Tasks (COMMENTED OUT)
    # ax1 = axes[0]
    # all_medians = [r['all']['median'] for r in results]
    # all_p90s = [r['all']['p90'] for r in results]
    # all_p99s = [r['all']['p99'] for r in results]
    # 
    # bars1 = ax1.bar(x - width, all_medians, width, label='Median', color=colors['Median'], edgecolor='black', linewidth=1)
    # bars2 = ax1.bar(x, all_p90s, width, label='P90', color=colors['P90'], edgecolor='black', linewidth=1)
    # bars3 = ax1.bar(x + width, all_p99s, width, label='P99', color=colors['P99'], edgecolor='black', linewidth=1)
    # 
    # ax1.set_xlabel('Scheduling Algorithm', fontsize=12, fontweight='bold')
    # ax1.set_ylabel('Response Time (ms)', fontsize=12, fontweight='bold')
    # ax1.set_title('All Tasks - Response Time Statistics', fontsize=14, fontweight='bold', pad=10)
    # ax1.set_xticks(x)
    # ax1.set_xticklabels(algorithms)
    # ax1.legend(loc='upper left')
    # ax1.grid(True, alpha=0.3, axis='y', linestyle='--')
    
    # Plot 1: Short Tasks (now first row)
    ax2 = axes[0]
    short_medians = [r['short']['median'] if r['short'] else 0 for r in results]
    short_p90s = [r['short']['p90'] if r['short'] else 0 for r in results]
    short_p99s = [r['short']['p99'] if r['short'] else 0 for r in results]
    
    bars4 = ax2.bar(x - width, short_medians, width, label='Median', color=colors['Median'], edgecolor='black', linewidth=1)
    bars5 = ax2.bar(x, short_p90s, width, label='P90', color=colors['P90'], edgecolor='black', linewidth=1)
    bars6 = ax2.bar(x + width, short_p99s, width, label='P99', color=colors['P99'], edgecolor='black', linewidth=1)
    
    ax2.set_xlabel('Scheduling Algorithm', fontsize=12, fontweight='bold')
    ax2.set_ylabel('Response Time (ms)', fontsize=12, fontweight='bold')
    ax2.set_title('Short Tasks - Response Time Statistics', fontsize=14, fontweight='bold', pad=10)
    ax2.set_xticks(x)
    ax2.set_xticklabels(algorithms)
    ax2.legend(loc='upper left')
    ax2.grid(True, alpha=0.3, axis='y', linestyle='--')
    
    # Plot 2: Long Tasks (now second row)
    ax3 = axes[1]
    long_medians = [r['long']['median'] if r['long'] else 0 for r in results]
    long_p90s = [r['long']['p90'] if r['long'] else 0 for r in results]
    long_p99s = [r['long']['p99'] if r['long'] else 0 for r in results]
    
    bars7 = ax3.bar(x - width, long_medians, width, label='Median', color=colors['Median'], edgecolor='black', linewidth=1)
    bars8 = ax3.bar(x, long_p90s, width, label='P90', color=colors['P90'], edgecolor='black', linewidth=1)
    bars9 = ax3.bar(x + width, long_p99s, width, label='P99', color=colors['P99'], edgecolor='black', linewidth=1)
    
    ax3.set_xlabel('Scheduling Algorithm', fontsize=12, fontweight='bold')
    ax3.set_ylabel('Response Time (ms)', fontsize=12, fontweight='bold')
    ax3.set_title('Long Tasks - Response Time Statistics', fontsize=14, fontweight='bold', pad=10)
    ax3.set_xticks(x)
    ax3.set_xticklabels(algorithms)
    ax3.legend(loc='upper left')
    ax3.grid(True, alpha=0.3, axis='y', linestyle='--')
    
    # Improve layout
    plt.tight_layout()
    
    # Save the plot
    output_file = 'algorithm_comparison.png'
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"\n✓ Comparison plot saved as: {output_file}")
    
    # Show the plot
    plt.show()
    
    # Print summary statistics
    print("\n" + "="*80)
    print("RESPONSE TIME COMPARISON SUMMARY")
    print("="*80)
    print(f"{'Algorithm':<15} {'Metric':<10} {'All Tasks':<15} {'Short Tasks':<15} {'Long Tasks':<15}")
    print("-"*80)
    
    for result in results:
        algo = result['algorithm']
        print(f"{algo:<15} {'Median':<10} {result['all']['median']:>12.2f} ms  ", end='')
        if result['short']:
            print(f"{result['short']['median']:>12.2f} ms  ", end='')
        else:
            print(f"{'N/A':>12}    ", end='')
        if result['long']:
            print(f"{result['long']['median']:>12.2f} ms")
        else:
            print(f"{'N/A':>12}")
        
        print(f"{'':<15} {'P90':<10} {result['all']['p90']:>12.2f} ms  ", end='')
        if result['short']:
            print(f"{result['short']['p90']:>12.2f} ms  ", end='')
        else:
            print(f"{'N/A':>12}    ", end='')
        if result['long']:
            print(f"{result['long']['p90']:>12.2f} ms")
        else:
            print(f"{'N/A':>12}")
        
        print(f"{'':<15} {'P99':<10} {result['all']['p99']:>12.2f} ms  ", end='')
        if result['short']:
            print(f"{result['short']['p99']:>12.2f} ms  ", end='')
        else:
            print(f"{'N/A':>12}    ", end='')
        if result['long']:
            print(f"{result['long']['p99']:>12.2f} ms")
        else:
            print(f"{'N/A':>12}")
        print("-"*80)
    
    print("="*80 + "\n")


def main():
    if len(sys.argv) < 2:
        print("Usage: python plot_results.py <csv_file1> <csv_file2> ...")
        print("Example: python plot_results.py results/fifo_results_*.csv results/sjf_results_*.csv")
        print("\nNote: Files should follow the pattern 'algo_results_timestamp.csv'")
        sys.exit(1)
    
    csv_files = sys.argv[1:]
    
    print(f"Loading data from {len(csv_files)} CSV file(s)...")
    plot_response_time_comparison(csv_files)


if __name__ == "__main__":
    main()


