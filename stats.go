package main

import "fmt"

// ContainerStats holds computed stats for display.
type ContainerStats struct {
	Name      string
	ID        string
	CPUPerc   float64
	MemUsage  uint64
	MemLimit  uint64
	MemPerc   float64
	NetRx     uint64
	NetTx     uint64
}

// CalculateStats computes human-readable stats from the Docker API response.
func CalculateStats(name, id string, stats *StatsResponse) ContainerStats {
	cs := ContainerStats{
		Name: name,
		ID:   id,
	}

	// CPU percentage
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta > 0 && cpuDelta > 0 {
		cpuCount := stats.CPUStats.OnlineCPUs
		if cpuCount == 0 {
			cpuCount = 1
		}
		cs.CPUPerc = (cpuDelta / systemDelta) * float64(cpuCount) * 100.0
	}

	// Memory
	cs.MemUsage = stats.MemoryStats.Usage
	cs.MemLimit = stats.MemoryStats.Limit
	if cs.MemLimit > 0 {
		cs.MemPerc = float64(cs.MemUsage) / float64(cs.MemLimit) * 100.0
	}

	// Network I/O (aggregate all interfaces)
	for _, netStats := range stats.Networks {
		cs.NetRx += netStats.RxBytes
		cs.NetTx += netStats.TxBytes
	}

	return cs
}

// FormatBytes converts bytes to a human-readable string.
func FormatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
