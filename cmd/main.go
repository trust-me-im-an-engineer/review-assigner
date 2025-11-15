package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"review-assigner/internal/config"
	"review-assigner/internal/rest"
	"review-assigner/internal/service"
	"review-assigner/internal/storage/postgres"
	"review-assigner/internal/validator"
)

func main() {
	// Configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Logging
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})
	slog.SetDefault(slog.New(jsonHandler))
	slog.Info("set json logging to stdout", "level", cfg.LogLevel)

	// Initializations
	storage, err := postgres.New(context.Background(), *cfg.DB)
	if err != nil {
		slog.Error("failed to initialize postgres storage", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("closing storage pool...")
		storage.Close()
		slog.Info("storage closed")
	}()
	serviceInst := service.NewService(storage)
	validatorInst := validator.New()
	router := rest.NewRouter(serviceInst, validatorInst)
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// Start Server in a Goroutine

	// Signal Handling Channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine so the main function can listen on stopCh
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server running", "address", cfg.Address)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Block until signal or error
	select {
	case err := <-serverErrors:
		slog.Error("fatal error while serving http", "error", err)
		os.Exit(1)
	case <-stopChan:
		slog.Info("received stop signal, initiating graceful shutdown...")
	}

	// Graceful Shutdown Context
	// Create a context with a timeout for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Shut Down HTTP Server
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("http server shutdown failed", "error", err)
	} else {
		slog.Info("http server gracefully stopped")
	}

	slog.Info("application stopped gracefully")
}
