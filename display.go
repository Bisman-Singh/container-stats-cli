package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	clearScreen = "\033[H\033[2J"
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// DisplayStats renders the container stats table to the terminal.
func DisplayStats(stats []ContainerStats) {
	fmt.Print(clearScreen)

	fmt.Printf("%s%sContainer Stats%s  [%s]\n\n",
		colorBold, colorCyan, colorReset, time.Now().Format("15:04:05"))

	if len(stats) == 0 {
		fmt.Println("No running containers found.")
		return
	}

	// Header
	fmt.Printf("%s%-20s %-14s %-8s %-22s %-8s %-14s %-14s%s\n",
		colorBold,
		"NAME", "CONTAINER ID", "CPU %",
		"MEM USAGE / LIMIT", "MEM %",
		"NET RX", "NET TX",
		colorReset,
	)
	fmt.Println(strings.Repeat("-", 105))

	for _, s := range stats {
		name := truncateName(s.Name, 19)
		shortID := s.ID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		cpuColor := colorGreen
		if s.CPUPerc > 80 {
			cpuColor = colorRed
		} else if s.CPUPerc > 50 {
			cpuColor = colorYellow
		}

		memColor := colorGreen
		if s.MemPerc > 80 {
			memColor = colorRed
		} else if s.MemPerc > 50 {
			memColor = colorYellow
		}

		memUsage := fmt.Sprintf("%s / %s", FormatBytes(s.MemUsage), FormatBytes(s.MemLimit))

		fmt.Printf("%-20s %-14s %s%-8.2f%s %-22s %s%-8.2f%s %-14s %-14s\n",
			name,
			shortID,
			cpuColor, s.CPUPerc, colorReset,
			memUsage,
			memColor, s.MemPerc, colorReset,
			FormatBytes(s.NetRx),
			FormatBytes(s.NetTx),
		)
	}

	fmt.Printf("\n%d container(s) running\n", len(stats))
}

func truncateName(name string, maxLen int) string {
	// Docker container names start with "/" from the API
	name = strings.TrimPrefix(name, "/")
	if len(name) > maxLen {
		return name[:maxLen-3] + "..."
	}
	return name
}
