package main

import (
	"context"

	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/config"
	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/repository"
	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/server"
	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/service"
	"github.com/jackc/pgx/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		return
	}

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return
	}
	defer repo.Close()

	db, err := pgx.Connect(context.Background(), cfg.Database.DSN)
	if err != nil {
		return
	}
	defer db.Close(context.Background())

	connectorService := service.NewConnectorService(cfg, repo)
	server := server.NewServer(connectorService)

	if err := server.Start(); err != nil {
		return
	}
	defer server.Stop()
}
