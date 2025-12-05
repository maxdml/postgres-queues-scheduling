#!/usr/bin/env python3
"""
Plot response time comparison for DBOS queue scheduling results.

Compares average response times across different scheduling algorithms.
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


def plot_response_time_comparison(csv_files):
    """Generate bar chart comparing average response times across algorithms."""
    
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
        
        # Calculate average response time
        avg_response_time = df['response_time_ms'].mean()
        
        results.append({
            'algorithm': algo_name,
            'avg_response_time': avg_response_time,
            'file': csv_file
        })
        
        print(f"Loaded {csv_file}: {algo_name} - Avg Response Time: {avg_response_time:.2f} ms")
    
    if not results:
        print("Error: No valid CSV files found!")
        sys.exit(1)
    
    # Sort results by algorithm name for consistent display
    results = sorted(results, key=lambda x: x['algorithm'])
    
    # Extract data for plotting
    algorithms = [r['algorithm'] for r in results]
    response_times = [r['avg_response_time'] for r in results]
    
    # Create the bar chart
    fig, ax = plt.subplots(figsize=(12, 7))
    
    # Create bars with different colors
    colors = plt.cm.Set3(np.linspace(0, 1, len(algorithms)))
    bars = ax.bar(algorithms, response_times, color=colors, edgecolor='black', linewidth=1.5, alpha=0.8)
    
    # Customize the plot
    ax.set_xlabel('Scheduling Algorithm', fontsize=14, fontweight='bold')
    ax.set_ylabel('Average Response Time (ms)', fontsize=14, fontweight='bold')
    ax.set_title('Average Response Time Comparison Across Scheduling Algorithms', 
                 fontsize=16, fontweight='bold', pad=20)
    ax.grid(True, alpha=0.3, axis='y', linestyle='--')
    
    # Add value labels on top of each bar
    for bar, value in zip(bars, response_times):
        height = bar.get_height()
        ax.text(bar.get_x() + bar.get_width()/2., height,
                f'{value:.1f} ms',
                ha='center', va='bottom', fontsize=11, fontweight='bold')
    
    # Improve layout
    plt.tight_layout()
    
    # Save the plot
    output_file = 'algorithm_comparison.png'
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    print(f"\n✓ Comparison plot saved as: {output_file}")
    
    # Show the plot
    plt.show()
    
    # Print summary statistics
    print("\n" + "="*70)
    print("RESPONSE TIME COMPARISON SUMMARY")
    print("="*70)
    for result in results:
        print(f"{result['algorithm']:15s}: {result['avg_response_time']:8.2f} ms")
    print("="*70)
    
    # Find best and worst
    best = min(results, key=lambda x: x['avg_response_time'])
    worst = max(results, key=lambda x: x['avg_response_time'])
    
    print(f"\nBest Algorithm:  {best['algorithm']} ({best['avg_response_time']:.2f} ms)")
    print(f"Worst Algorithm: {worst['algorithm']} ({worst['avg_response_time']:.2f} ms)")
    
    if len(results) > 1:
        improvement = ((worst['avg_response_time'] - best['avg_response_time']) / 
                      worst['avg_response_time'] * 100)
        print(f"Improvement:     {improvement:.1f}% reduction in response time")
    
    print("="*70 + "\n")


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


