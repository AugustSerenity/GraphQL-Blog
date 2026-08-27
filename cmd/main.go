package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/AugustSerenity/GraphQL-Blog/internal/auth"
	"github.com/AugustSerenity/GraphQL-Blog/internal/config"
	"github.com/AugustSerenity/GraphQL-Blog/internal/graph"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository/memory"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository/postgres"
	"github.com/AugustSerenity/GraphQL-Blog/internal/service"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	slog.SetDefault(logger)

	cfg := config.Load()

	repo, closeRepo := createRepository(cfg.RepositoryType, cfg.DatabaseURL, logger)
	defer closeRepo()

	svc := service.NewService(repo)

	resolver := &graph.Resolver{
		Service: svc,
	}

	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolver,
			},
		),
	)

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 5 * time.Second,
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	http.Handle("/query", auth.Middleware(srv))

	logger.Info("server started", "address", ":8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func createRepository(repositoryType, databaseURL string, logger *slog.Logger) (repository.Repository, func()) {
	switch repositoryType {
	case "", "inmemory":
		logger.Info("using inmemory repository")

		repo := memory.New()

		return repo, func() {}

	case "postgres":
		if databaseURL == "" {
			logger.Error("DATABASE_URL is required when REPOSITORY_TYPE=postgres")
			os.Exit(1)
		}

		db, err := postgres.InitDB(databaseURL, logger)
		if err != nil {
			logger.Error("failed to initialize postgres", "error", err)
			os.Exit(1)
		}

		repo := postgres.New(db)

		return repo, func() {
			postgres.CloseDB(db, logger)
		}

	default:
		logger.Error(
			"unknown repository type",
			"repository",
			repositoryType,
			"allowed",
			"memory, postgres",
		)
		os.Exit(1)

		return nil, nil
	}
}
