package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/AugustSerenity/GraphQL-Blog/internal/graph"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository/memory"
	"github.com/AugustSerenity/GraphQL-Blog/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	// cfg := config.Load()

	// db, err := postgres.InitDB(cfg.DatabaseURL, logger)
	// if err != nil {
	// 	logger.Error("failed to initialize database", "error", err)
	// 	os.Exit(1)
	// }
	// defer postgres.CloseDB(db, logger)

	repo := memory.New()
	svc := service.NewService(repo)

	resolver := &graph.Resolver{
		Service: svc,
	}

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolver,
			},
		),
	)

	http.Handle("/query", srv)

	logger.Info("server started", "address", ":8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
