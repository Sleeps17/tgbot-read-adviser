package main

import (
	"log"
	"log/slog"
	"os"

	tgClient "tgbot-read-adviser/internal/clients/telegram"
	"tgbot-read-adviser/internal/config"
	"tgbot-read-adviser/internal/consumer/event_consumer"
	"tgbot-read-adviser/internal/events/telegram"
	"tgbot-read-adviser/internal/storage/sqlite"
)

func main() {

	cfg := config.MustLoad()
	logger := setupLogger(cfg)
	logger.Info("the logger is successfully configured")

	s, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to create storage: ", slog.String("[ERR]", err.Error()))
		os.Exit(0)
	}

	eventsProcessor := telegram.New(
		tgClient.New(cfg.TgBotHost, cfg.TgBotToken),
		s,
	)
	logger.Info("the event_processor is successfully configured")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, cfg.BatchSize)
	logger.Info("the consumer is successfully configured")

	logger.Info("service started")
	if err := consumer.Start(); err != nil {
		logger.Error("service is stopped", slog.String("[ERR]", err.Error()))
		os.Exit(0)
	}
}

const (
	LocalEnv = "local"
	DevEnv   = "dev"
	ProdEnv  = "prod"
)

func setupLogger(cfg *config.Config) *slog.Logger {

	logs, err := os.OpenFile(cfg.LogsPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error with open logs file: %v", err)
	}

	var logger *slog.Logger

	switch cfg.Env {

	case LocalEnv:
		logger = slog.New(
			slog.NewTextHandler(logs, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}),
		)

	case DevEnv:
		logger = slog.New(
			slog.NewJSONHandler(logs, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case ProdEnv:
		logger = slog.New(
			slog.NewJSONHandler(logs, &slog.HandlerOptions{Level: slog.LevelError}),
		)

	default:
		logger = slog.New(
			slog.NewTextHandler(logs, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
