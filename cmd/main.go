package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/config"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/client/jira"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/repository/postgres"
	mysync "github.com/Go-Yadro-Group-1/Jira-Connector/internal/service/sync"
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
	svc := mysync.NewService(jiraClient, repo)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	jql := `project = "` + *projectKey + `"`

	errChan := make(chan error, 100)
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			childCtx, childCancel := context.WithTimeout(ctx, time.Second*30)
			defer childCancel()

			if err := svc.RunWorkerPool(childCtx, jql, 50); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()
}
