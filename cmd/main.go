package main

import (
	"github.com/AugustSerenity/GraphQL-Blog/internal/config"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository"
)

func main() {
	cfg := config.Load()

	db := repository.InitDB(cfg.DBPath)
	defer repository.CloseDB(db)
}
