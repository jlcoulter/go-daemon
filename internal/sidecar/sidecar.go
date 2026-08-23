package sidecar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

// Sidecar provides an HTTP server for health checks and metrics
// alongside the main daemon. Runs on a separate port.
type Sidecar struct {
	server *http.Server
	port   int
	ready  atomic.Bool
}

func New(port int) *Sidecar {
	s := &Sidecar{port: port}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/readyz", s.readyz)
	mux.HandleFunc("/stats", s.stats)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return s
}

// SetReady marks the daemon as ready to accept traffic.
func (s *Sidecar) SetReady(ready bool) {
	s.ready.Store(ready)
}

func (s *Sidecar) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Sidecar) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Sidecar) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Sidecar) readyz(w http.ResponseWriter, r *http.Request) {
	if s.ready.Load() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "not ready"})
	}
}

func (s *Sidecar) stats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: Expose daemon metrics (active connections, uptime, etc.)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}