package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		panic(errors.Wrap(err, "failed to load config"))
	}

	logger := buildLogger(cfg)
	slog.SetDefault(logger)

	logger.Info("Starting app", "configs", cfg.LogFields())

	appCtx, cancel := context.WithCancel(context.Background())

	application, err := NewApplication(appCtx, cfg, logger)
	if err != nil {
		logger.Error("failed to create application", "error", err)
		return
	}

	if err = application.Run(appCtx); err != nil {
		logger.Error("failed to run application:", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	logger.Info("Gracefully shutting down...")

	cancel()

	if err = application.Stop(); err != nil {
		logger.Error("failed to stop application: %v", err)
		return
	}
}

func buildLogger(cfg *Config) *slog.Logger {
	logLevel := slog.LevelInfo

	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}
