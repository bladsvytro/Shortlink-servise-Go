package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	zapLogger, err := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// Create application
	application, err := app.New(cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("Failed to create application", zap.Error(err))
	}

	// Start application
	if err := application.Start(); err != nil {
		zapLogger.Fatal("Failed to start application", zap.Error(err))
	}

	zapLogger.Info("Application started",
		zap.Int("port", cfg.Server.Port),
		zap.String("environment", cfg.Server.Env),
	)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Stop(ctx); err != nil {
		zapLogger.Error("Failed to stop application gracefully", zap.Error(err))
	} else {
		zapLogger.Info("Server stopped gracefully")
	}
}
