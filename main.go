package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	socketPath := flag.String("socket", "/var/run/docker.sock", "path to Docker socket")
	interval := flag.Duration("interval", 2*time.Second, "refresh interval")
	once := flag.Bool("once", false, "fetch stats once and exit")
	flag.Parse()

	// Check if socket exists
	if _, err := os.Stat(*socketPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Docker socket not found at %s\n", *socketPath)
		fmt.Fprintln(os.Stderr, "Make sure Docker is running or specify the socket path with -socket")
		os.Exit(1)
	}

	client := NewDockerClient(*socketPath)

	if *once {
		stats, err := fetchAllStats(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		DisplayStats(stats)
		return
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	// Initial fetch
	stats, err := fetchAllStats(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	DisplayStats(stats)

	for {
		select {
		case <-ticker.C:
			stats, err := fetchAllStats(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
				continue
			}
			DisplayStats(stats)
		case <-sigChan:
			fmt.Println("\nShutting down...")
			os.Exit(0)
		}
	}
}

func fetchAllStats(client *DockerClient) ([]ContainerStats, error) {
	containers, err := client.ListContainers()
	if err != nil {
		return nil, err
	}

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []ContainerStats
	)

	for _, container := range containers {
		wg.Add(1)
		go func(c Container) {
			defer wg.Done()

			stats, err := client.GetStats(c.ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to get stats for %s: %v\n", c.ID[:12], err)
				return
			}

			name := c.ID[:12]
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}

			cs := CalculateStats(name, c.ID, stats)

			mu.Lock()
			results = append(results, cs)
			mu.Unlock()
		}(container)
	}

	wg.Wait()
	return results, nil
}
