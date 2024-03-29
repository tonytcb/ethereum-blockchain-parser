package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.Default()

	logger.Info("Starting app test...")

	cfg, err := loadConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return
	}

	logLevel := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	slog.SetDefault(logger)

	logger.Info("Application configurations", "data", cfg.LogFields())

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

	logger.Info("gracefully shutting down...")

	cancel()

	if err = application.Stop(); err != nil {
		logger.Error("failed to stop application: %v", err)
		return
	}
}
