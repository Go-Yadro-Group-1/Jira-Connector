package repository

import (
	"github.com/Go-Yadro-Group-1/Jira-Analyzer/internal/config"
)

type Repository struct {
	cfg *config.Config
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	return &Repository{
		cfg: cfg,
	}, nil
}

func (r *Repository) Close() error {
	return nil
}
