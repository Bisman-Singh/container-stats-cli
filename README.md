# Container Stats CLI

A terminal-based container stats viewer written in Go. Reads container statistics directly from the Docker socket API and displays CPU, memory, and network I/O in a refreshing table format.

## Features

- Reads stats via Docker socket (`/var/run/docker.sock`)
- Displays container name, CPU%, memory usage/limit, memory%, network I/O
- Color-coded resource usage (green/yellow/red thresholds)
- Auto-refreshes every 2 seconds (configurable)
- Concurrent stats fetching for all containers
- Single-run mode for scripting
- No external dependencies, uses only the Go standard library

## Build

```bash
go build -o container-stats-cli .
```

## Usage

```bash
# Run with defaults (refresh every 2s)
./container-stats-cli

# Custom refresh interval
./container-stats-cli -interval 5s

# Fetch stats once and exit
./container-stats-cli -once

# Use a custom Docker socket path
./container-stats-cli -socket /path/to/docker.sock
```

### Flags

| Flag        | Default                    | Description                   |
|-------------|----------------------------|-------------------------------|
| `-socket`   | `/var/run/docker.sock`     | Path to Docker socket         |
| `-interval` | `2s`                       | Refresh interval              |
| `-once`     | `false`                    | Fetch stats once and exit     |

### Example Output

```
Container Stats  [14:23:01]

NAME                 CONTAINER ID   CPU %    MEM USAGE / LIMIT      MEM %    NET RX         NET TX
---------------------------------------------------------------------------------------------------------
web-app              a1b2c3d4e5f6   12.45    256.00 MB / 2.00 GB    12.50    45.23 MB       12.34 MB
redis-cache          f6e5d4c3b2a1   2.30     64.00 MB / 512.00 MB   12.50    1.23 MB        456.00 KB
postgres-db          1a2b3c4d5e6f   8.75     512.00 MB / 4.00 GB    12.50    23.45 MB       67.89 MB

3 container(s) running
```

## Requirements

- Docker must be running
- Access to the Docker socket (typically requires root or docker group membership)


