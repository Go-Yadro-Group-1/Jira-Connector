package sync

import (
	"context"
	"fmt"

	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/client/jira"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/repository/postgres"
)

type SyncService struct {
	jiraClient *jira.JiraClient
	repo       *postgres.PostgresRepository
}

type ResultWithID struct {
	TaskID string
	Result[jira.Issue]
}

func NewService(jiraClient *jira.JiraClient, repo *postgres.PostgresRepository) *SyncService {
	return &SyncService{jiraClient: jiraClient, repo: repo}
}

func (s *SyncService) RunWorkerPool(ctx context.Context, jql string, maxWorkers int) error {
	searchResp, err := s.jiraClient.SearchIssues(ctx, jql)
	if err != nil {
		return fmt.Errorf("failed to search issues: %w", err)
	}

	if len(searchResp.Issues) == 0 {
		return fmt.Errorf("no issues found matching JQL query")
	}

	fmt.Printf("Found %d issues to process\n", len(searchResp.Issues))

	pool := New[jira.Issue](maxWorkers)
	pool.Start(ctx)
	// defer pool.Stop()

	taskCount := len(searchResp.Issues)
	resultChan := make(chan ResultWithID, taskCount)

	go func() {
		for i := range searchResp.Issues {
			issue := searchResp.Issues[i]
			taskID := issue.Key

			pool.Submit(func(ctx context.Context) (jira.Issue, error) {
				issuePtr, err := s.jiraClient.GetIssue(ctx, issue.Key)

				if err != nil {
					return jira.Issue{}, err
				}

				return *issuePtr, nil
			})

			resultChan <- ResultWithID{TaskID: taskID}
		}

		close(pool.tasks)
	}()

	var errs []error
	completed := 0
	taskIDs := make(map[int]string)

	for i := range searchResp.Issues {
		taskIDs[i] = searchResp.Issues[i].Key
	}

	for completed < taskCount {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok := <-pool.Results():
			if !ok {
				break
			}
			completed++

			if res.Err != nil {
				errs = append(errs, fmt.Errorf("task %s (processed by %s) failed: %w",
					taskIDs[completed-1], res.ID, res.Err))
				continue
			}

			fmt.Printf("[%s] Successfully processed %s: %s (Status: %s)\n",
				res.ID,
				res.Value.Key,
				res.Value.Fields.Summary,
				res.Value.Fields.Status.Name)
		}
	}

	if len(errs) > 0 {
		fmt.Printf("\n%d errors occurred during processing:\n", len(errs))
		for i, err := range errs {
			fmt.Printf("  %d. %v\n", i+1, err)
		}
		return fmt.Errorf("%d tasks failed", len(errs))
	}

	fmt.Printf("\nSuccessfully processed %d issues\n", taskCount)
	return nil
}
