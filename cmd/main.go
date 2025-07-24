package main

import (
	"context"
	"os"
	"os/signal"
	"rent-application/configs"
	"rent-application/internal/logger"
	"rent-application/internal/server"
	"syscall"
	"time"
)

func main() {

	// Load config
	filename := "internal/config/config.yaml"
	if err := configs.LoadConfig(filename); err != nil {
		panic(err)
	}

	// Initialize logger (zerolog)
	logger := logger.New("development")
	defer logger.Sync()

	// Create server
	srv, err := server.New(configs.Cfg, &logger.Logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create server")
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error().Err(err).Msg("Server error")
		}
	}()

	<-quit
	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Forced shutdown")
	}

	logger.Info().Msg("Server stopped")
}
