package app

import (
	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/config"
)

type App struct {
	cfg config.Config
}

func New(cfg config.Config) (*App, error) {
	return &App{
		cfg: cfg,
	}, nil
}

func (a *App) Run() <-chan error {
	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		errChan <- a.run()
	}()
	return errChan
}

func (a *App) Close() error {
	return nil
}

func (a *App) run() error {
	return nil
}
