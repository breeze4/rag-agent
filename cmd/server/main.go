package main

import (
	"log/slog"
	"os"
	"rag-therapist/internal/config"
)

func main() {
	// Set up structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("RAG Therapist server starting...")
	
	// Load configuration
	cfg := config.Load()
	
	slog.Info("Server initialized", "port", cfg.Port)
	
	// TODO: Initialize server
}