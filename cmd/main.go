package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/AugustSerenity/GraphQL-Blog/internal/graph"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository/memory"
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

	repo := memory.New()
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

	http.Handle("/query", srv)

	logger.Info("server started", "address", ":8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
