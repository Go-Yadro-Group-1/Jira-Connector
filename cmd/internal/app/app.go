package app

import (
	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/config"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/broker/consumer"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/broker/publisher"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/service/sync"
)

type App struct {
	cfg       config.Config
	consumer  *consumer.Consumer
	publisher *publisher.Publisher
	syncer    *sync.SyncService
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
