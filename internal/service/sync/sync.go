package sync

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/client/jira"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/repository/postgres"
)

var (
	ErrNoIssuesFound = errors.New("no issues found matching JQL query")
	ErrTasksFailed   = errors.New("one or more tasks failed")
)

type Service struct {
	jiraClient *jira.Client
	repo       *postgres.PostgresRepository
}

func NewService(jiraClient *jira.Client, repo *postgres.PostgresRepository) *Service {
	return &Service{jiraClient: jiraClient, repo: repo}
}

func (s *Service) RunWorkerPool(ctx context.Context, jql string, maxWorkers int) error {
	searchResp, err := s.jiraClient.SearchIssues(ctx, jql)
	if err != nil {
		return fmt.Errorf("failed to search issues: %w", err)
	}

	if len(searchResp.Issues) == 0 {
		return ErrNoIssuesFound
	}

	log.Printf("Found %d issues to process\n", len(searchResp.Issues))

	pool := New[jira.Issue](maxWorkers)
	pool.Start(ctx)
	// defer pool.Stop()

	taskCount := len(searchResp.Issues)

	go func() {
		defer close(pool.tasks)
		for _, issue := range searchResp.Issues {
			key := issue.Key

			pool.Submit(func(ctx context.Context) (jira.Issue, error) {
				issuePtr, err := s.jiraClient.GetIssue(ctx, key)
				if err != nil {
					return jira.Issue{}, fmt.Errorf("failed to get issue %q: %w", key, err)
				}

				return *issuePtr, nil
			})
		}
	}()

	return s.processResults(ctx, pool, taskCount)
}

func (s *Service) processResults(
	ctx context.Context,
	pool *WorkerPool[jira.Issue],
	taskCount int,
) error {
	var errs []error
	completed := 0

	for completed < taskCount {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case res, ok := <-pool.Results():
			if !ok {
				break
			}

			completed++

			if res.Err != nil {
				errs = append(errs, res.Err)

				continue
			}

			issue := res.Value

			log.Printf("[%s] Successfully processed %s: %s (Status: %s)\n",
				res.ID,
				issue.Key,
				issue.Fields.Summary,
				issue.Fields.Status.Name)

			// s.repo.SaveIssue(issue)
		}
	}

	if len(errs) > 0 {
		log.Printf("\n%d errors occurred during processing:\n", len(errs))

		for i, err := range errs {
			log.Printf("  %d. %v\n", i+1, err)
		}

		return ErrTasksFailed
	}

	log.Printf("\nSuccessfully processed %d issues\n", taskCount)

	return nil
}
