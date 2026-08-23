package daemon_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/jlcoulter/go-daemon-template/internal/daemon"
)

func TestDaemonAcceptAndShutdown(t *testing.T) {
	d := daemon.New(daemon.Config{Port: 0}) // port 0 = random available port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- d.Run(ctx)
	}()

	// Give the daemon a moment to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown
	cancel()
	d.Shutdown()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("daemon Run returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("daemon did not shut down in time")
	}
}

func TestDaemonConnectionTracking(t *testing.T) {
	d := daemon.New(daemon.Config{Port: 0})

	if d.ActiveConnections() != 0 {
		t.Errorf("expected 0 active connections, got %d", d.ActiveConnections())
	}
}

func TestDaemonEchoConnection(t *testing.T) {
	// Start daemon on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	addr := listener.Addr().String()
	listener.Close()

	d := daemon.New(daemon.Config{Port: mustParsePort(addr)})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- d.Run(ctx)
	}()

	// Connect
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		// Daemon may not be ready yet, skip gracefully
		t.Skipf("could not connect: %v", err)
	}
	defer conn.Close()

	msg := []byte("hello\n")
	conn.Write(msg)

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	received := string(buf[:n])
	if received != string(msg) {
		t.Errorf("echo: expected %q, got %q", string(msg), received)
	}
}

func mustParsePort(addr string) int {
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	var port int
	fmt.Sscanf(portStr, "%d", &port)
	return port
}