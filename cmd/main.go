package main

import (
	"flag"
	"log"

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

	logrus.Info("Starting Jira sync")

	err = svc.SyncProject(*projectKey)
	if err != nil {
		logrus.Errorf("Sync failed: %v", err)
	} else {
		logrus.Info("Sync completed successfully!")
	}
}
