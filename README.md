# Go Daemon Template

A GitHub template for long-running Go services — SSH servers, SMTP relays, gRPC proxies, TCP agents, watchers, and anything that stays alive rather than handling HTTP requests.

## Why not go-api?

`go-api` assumes request/response HTTP with a chi router. Daemons are different:
- They listen on non-HTTP protocols (SSH, SMTP, gRPC, TCP)
- They manage long-lived connections, not short requests
- They need graceful shutdown that drains active connections
- They often run as system services with PID management
- They may expose a sidecar HTTP endpoint for health/metrics, but it's secondary

## Features

- **Graceful shutdown** — SIGINT/SIGTERM drains active connections before exit
- **Sidecar HTTP server** — optional `/healthz` and `/metrics` on a separate port
- **Connection tracking** — track and drain active connections on shutdown
- **Structured logging** with `log/slog`
- **Configuration** via environment variables with viper
- **Systemd unit file** included
- **Docker** multi-stage build
- **CI** via GitHub Actions
- **Release** via GoReleaser

## Usage

1. Click **"Use this template"** on GitHub to create a new repo
2. Run the setup script:
   ```sh
   ./setup.sh mydaemon github.com/you/mydaemon
   ```
3. Replace `internal/daemon/daemon.go` with your protocol logic
4. Add your connection handlers

## Project Structure

```
.
├── cmd/
│   └── daemon/
│       └── main.go          # Entry point, signal handling, graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go        # Viper config (env vars + defaults)
│   ├── daemon/
│   │   └── daemon.go         # Core daemon logic (replace this)
│   └── sidecar/
│       └── sidecar.go        # Optional HTTP health/metrics server
├── .github/
│   └── workflows/
│       └── ci.yml
├── contrib/
│   └── mydaemon.service      # Systemd unit file
├── .goreleaser.yml
├── Dockerfile
├── Makefile
├── go.mod
├── setup.sh
└── README.md
```

## Quick Start

```sh
# Run locally
make run

# Run tests
make test

# Build binary
make build

# Build Docker image
make docker

# Install as systemd service
sudo cp contrib/mydaemon.service /etc/systemd/system/
sudo systemctl enable --now mydaemon
```

## Container Images

CI builds and pushes a container image to GHCR on every push to any branch.

```sh
# Pull the latest image
docker pull ghcr.io/<owner>/go-daemon-template:latest

# Pull a specific commit
docker pull ghcr.io/<owner>/go-daemon-template:<sha>

# Run
docker run -p 2222:2222 -p 8080:8080 ghcr.io/<owner>/go-daemon-template:latest
```

Replace `<owner>` with your GitHub username or org. Images are tagged with both `latest` and the commit SHA.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Sidecar HTTP port (health/metrics) |
| `DAEMON_PORT` | `2222` | Main daemon listen port |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |

## Make Targets

| Target | Description |
|--------|-------------|
| `make run` | Run the daemon locally |
| `make build` | Build the binary to `./bin/` |
| `make test` | Run tests |
| `make vet` | Run `go vet` |
| `make lint` | Run golangci-lint |
| `make docker` | Build Docker image |
| `make clean` | Remove build artifacts |

## License

MIT