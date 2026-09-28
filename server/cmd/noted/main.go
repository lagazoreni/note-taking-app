package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"noted.local/noted/internal/api"
	"noted.local/noted/internal/api/handlers"
	"noted.local/noted/internal/config"
	"noted.local/noted/internal/store/sqlite"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := config.FromEnvironment()
	if err != nil {
		return err
	}
	if len(args) > 0 {
		switch args[0] {
		case "healthcheck":
			return healthcheck(cfg)
		case "db-status":
			return databaseCommand(cfg, "status")
		case "db-check":
			return databaseCommand(cfg, "check")
		case "db-reset":
			return databaseCommand(cfg, "reset")
		}
	}
	logger := newLogger(cfg.LogLevel)
	ctx := context.Background()
	if err := os.MkdirAll(cfg.DataDir, 0o750); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	db, err := sqlite.OpenContext(ctx, cfg.DatabasePath())
	if err != nil {
		return err
	}
	defer db.Close()
	if err := sqlite.RunMigrations(ctx, db); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}
	ready := true
	var assets fs.FS
	if stat, err := os.Stat("webdist"); err == nil && stat.IsDir() {
		assets = os.DirFS("webdist")
	}
	handler := api.NewRouter(api.RouterOptions{DB: db, Version: cfg.Version, Ready: func() bool { return ready }, Assets: assets, API: handlers.NewAPI(db), MaxBodyBytes: cfg.ImportMaxBytes, Logger: logger})
	server := &http.Server{Addr: cfg.Addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 60 * time.Second, IdleTimeout: 120 * time.Second}
	serverErr := make(chan error, 1)
	go func() { logger.Info("server started", "addr", cfg.Addr); serverErr <- server.ListenAndServe() }()
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case <-signalCtx.Done():
		ready = false
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = sqlite.Optimize(shutdownCtx, db)
		return server.Shutdown(shutdownCtx)
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func newLogger(level string) *slog.Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}

func healthcheck(cfg config.Config) error {
	client := http.Client{Timeout: 2 * time.Second}
	address := cfg.Addr
	if _, port, splitErr := net.SplitHostPort(cfg.Addr); splitErr == nil {
		address = "127.0.0.1:" + port
	}
	resp, err := client.Get("http://" + address + "/healthz")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned %s", resp.Status)
	}
	return nil
}

func databaseCommand(cfg config.Config, command string) error {
	if command == "reset" {
		if filepath.Clean(cfg.DataDir) == "." || filepath.Clean(cfg.DataDir) == string(filepath.Separator) {
			return fmt.Errorf("refusing to reset unsafe data directory")
		}
		if err := os.Remove(cfg.DatabasePath()); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	db, err := sqlite.OpenContext(context.Background(), cfg.DatabasePath())
	if err != nil {
		return err
	}
	defer db.Close()
	if err := sqlite.RunMigrations(context.Background(), db); err != nil {
		return err
	}
	if command == "check" {
		if err := checkDatabase(context.Background(), db); err != nil {
			return err
		}
		fmt.Println("ok")
		return nil
	}
	var count int
	if err := db.QueryRow("SELECT count(*) FROM schema_migrations").Scan(&count); err != nil {
		return err
	}
	fmt.Printf("migrations=%d database=%s\n", count, cfg.DatabasePath())
	return nil
}
