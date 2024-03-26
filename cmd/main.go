package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tonytcb/ethereum-blockchain-parser/internal/app"
	"github.com/tonytcb/ethereum-blockchain-parser/internal/app/config"
)

func main() {
	logger := slog.Default()

	logger.Info("starting app...")

	cfg, err := config.Load("./config/default.env", "default.env")
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return
	}

	logger.Info("Application env vars", "data", cfg.LogFields())

	appCtx, cancel := context.WithCancel(context.Background())

	application, err := app.NewApplication(appCtx, cfg)
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
