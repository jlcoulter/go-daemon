package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/jlcoulter/go-daemon-template/internal/config"
	"github.com/jlcoulter/go-daemon-template/internal/daemon"
	"github.com/jlcoulter/go-daemon-template/internal/sidecar"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "go-daemon-template",
	Short: "A long-running Go daemon",
	Long:  "A long-running Go daemon template — replace this with your service description.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := slog.LevelInfo
		if verbose {
			level = slog.LevelDebug
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})))
	},
	RunE: runDaemon,
}

func runDaemon(cmd *cobra.Command, args []string) error {
	cfg := config.Load(cfgFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start sidecar HTTP server for health/metrics
	sidecarSrv := sidecar.New(cfg.SidecarPort)
	go func() {
		slog.Info("starting sidecar HTTP server", "port", cfg.SidecarPort)
		if err := sidecarSrv.ListenAndServe(); err != nil {
			slog.Error("sidecar server error", "error", err)
		}
	}()

	// Start the main daemon
	d := daemon.New(daemon.Config{Port: cfg.DaemonPort})
	go func() {
		slog.Info("starting daemon", "port", cfg.DaemonPort)
		if err := d.Run(ctx); err != nil {
			slog.Error("daemon error", "error", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("daemon running, press Ctrl+C to stop")

	sig := <-quit
	slog.Info("received signal, shutting down", "signal", sig)

	// Graceful shutdown: stop accepting new connections, drain active ones
	slog.Info("shutting down sidecar")
	if err := sidecarSrv.Shutdown(ctx); err != nil {
		slog.Error("sidecar shutdown error", "error", err)
	}

	slog.Info("shutting down daemon, draining connections")
	if err := d.Shutdown(); err != nil {
		slog.Error("daemon shutdown error", "error", err)
	}

	slog.Info("daemon stopped")
	return nil
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func parseInt(s string, def int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}