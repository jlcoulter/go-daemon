package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
)

// Daemon manages the main long-running service.
// Replace Run() with your protocol logic (SSH, SMTP, gRPC, TCP, etc.)
type Daemon struct {
	cfg     Config
	active  atomic.Int64
	connWg  sync.WaitGroup
	mu      sync.Mutex
	running bool
}

type Config struct {
	Port int
}

func New(cfg Config) *Daemon {
	return &Daemon{
		cfg: cfg,
	}
}

// Run starts the daemon listener. Blocks until context is cancelled or error.
func (d *Daemon) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", d.cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	defer listener.Close()

	d.running = true
	slog.Info("daemon listening", "addr", addr)

	// Accept loop
	connCh := make(chan net.Conn)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// Listener closed — shutdown
				close(connCh)
				return
			}
			connCh <- conn
		}
	}()

	for {
		select {
		case <-ctx.Done():
			slog.Info("daemon accept loop cancelled")
			listener.Close()
			return nil
		case conn, ok := <-connCh:
			if !ok {
				return nil
			}
			d.active.Add(1)
			d.connWg.Add(1)
			go d.handleConn(conn)
		}
	}
}

// handleConn processes a single connection.
// Replace this with your protocol logic.
func (d *Daemon) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		d.active.Add(-1)
		d.connWg.Done()
	}()

	slog.Info("new connection", "remote", conn.RemoteAddr())

	// TODO: Replace with your protocol handler
	// Example: SSH, SMTP, gRPC, custom protocol
	// For now, just read and echo
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		conn.Write(buf[:n])
	}
}

// Shutdown gracefully drains active connections.
func (d *Daemon) Shutdown() error {
	d.running = false
	slog.Info("waiting for active connections", "count", d.active.Load())
	d.connWg.Wait()
	return nil
}

// ActiveConnections returns the number of active connections.
func (d *Daemon) ActiveConnections() int64 {
	return d.active.Load()
}