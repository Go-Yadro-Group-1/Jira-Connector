package service

import (
	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/config"
	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/repository"
)

type Connector struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewConnectorService(cfg *config.Config, repo *repository.Repository) *Connector {
	return &Connector{
		cfg:  cfg,
		repo: repo,
	}
}
