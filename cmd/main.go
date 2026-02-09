package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/config"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/client/jira"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/repository/postgres"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/service/sync"
	"github.com/sirupsen/logrus"
)

var projectKey = flag.String("project", "", "Jira project key")

func main() {
	flag.Parse()

	if *projectKey == "" {
		log.Fatal("Usage: go run main.go --project=PROJ")
	}

	cfg, err := config.LoadDevConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	jiraClient := jira.New(cfg.Jira.BaseURL, cfg.Jira.Token)
	repo := postgres.NewRepository()
	svc := sync.NewService(jiraClient, repo)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logrus.Info("Received interrupt signal, shutting down...")
		cancel()
	}()

	jql := `project = "` + *projectKey + `"`

	if err := svc.RunWorkerPool(ctx, jql, 10); err != nil {
		logrus.WithError(err).Error("Sync failed")
		os.Exit(1)
	}

	logrus.Info("Sync completed successfully!")
}
