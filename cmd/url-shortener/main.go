package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	logger := setupLogger(cfg.Env)

	logger.Info("starting server", "address", cfg.HTTPServer.Address)

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		logger.Error("cannot create storage", sl.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveURL("http://google.com", "google")

	if err != nil {
		logger.Error("cannot save url", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("url saved", "id", id)

	id, err = storage.SaveURL("http://google.com", "google")

	if err != nil {
		logger.Error("cannot save url", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("url saved", "id", id)
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
