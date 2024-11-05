package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

// DockerClient communicates with the Docker daemon via the Unix socket.
type DockerClient struct {
	httpClient *http.Client
	socketPath string
}

// Container represents a Docker container from the list API.
type Container struct {
	ID    string   `json:"Id"`
	Names []string `json:"Names"`
	State string   `json:"State"`
}

// StatsResponse represents the Docker stats API response.
type StatsResponse struct {
	Read    string `json:"read"`
	PreRead string `json:"preread"`

	CPUStats    CPUStatsEntry    `json:"cpu_stats"`
	PreCPUStats CPUStatsEntry    `json:"precpu_stats"`
	MemoryStats MemoryStatsEntry `json:"memory_stats"`
	Networks    map[string]NetworkStatsEntry `json:"networks"`
}

// CPUStatsEntry holds CPU usage data.
type CPUStatsEntry struct {
	CPUUsage    CPUUsage    `json:"cpu_usage"`
	SystemUsage uint64      `json:"system_cpu_usage"`
	OnlineCPUs  int         `json:"online_cpus"`
}

// CPUUsage holds per-container CPU usage.
type CPUUsage struct {
	TotalUsage uint64 `json:"total_usage"`
}

// MemoryStatsEntry holds memory usage data.
type MemoryStatsEntry struct {
	Usage uint64 `json:"usage"`
	Limit uint64 `json:"limit"`
}

// NetworkStatsEntry holds network I/O data for a single interface.
type NetworkStatsEntry struct {
	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`
}

// NewDockerClient creates a new client connected to the Docker socket.
func NewDockerClient(socketPath string) *DockerClient {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}

	return &DockerClient{
		httpClient: &http.Client{Transport: transport},
		socketPath: socketPath,
	}
}

// ListContainers returns all running containers.
func (d *DockerClient) ListContainers() ([]Container, error) {
	resp, err := d.httpClient.Get("http://localhost/containers/json")
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("docker API error (status %d): %s", resp.StatusCode, string(body))
	}

	var containers []Container
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, fmt.Errorf("failed to parse container list: %w", err)
	}

	return containers, nil
}

// GetStats fetches stats for a single container (one-shot, not streaming).
func (d *DockerClient) GetStats(containerID string) (*StatsResponse, error) {
	url := fmt.Sprintf("http://localhost/containers/%s/stats?stream=false", containerID)
	resp, err := d.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats for %s: %w", containerID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("docker API error (status %d): %s", resp.StatusCode, string(body))
	}

	var stats StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	return &stats, nil
}
