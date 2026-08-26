package main

import (
	"log/slog"
	"os"

	"github.com/AugustSerenity/GraphQL-Blog/internal/config"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository/postgres"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := postgres.InitDB(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer postgres.CloseDB(db, logger)
}
