package sync

import (
	"fmt"
	"log"

	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/client/jira"
	"github.com/Go-Yadro-Group-1/Jira-Connector/internal/repository/postgres"
)

type Service struct {
	jiraClient *jira.JiraClient
	repo       *postgres.Repository
}

func NewService(jiraClient *jira.JiraClient, repo *postgres.Repository) *Service {
	return &Service{jiraClient: jiraClient, repo: repo}
}

func (s *Service) SyncProject(projectKey string) error {
	sr, err := s.jiraClient.SearchIssues(fmt.Sprintf(`project = "%s"`, projectKey))
	if err != nil {
		return err
	}

	for _, issue := range sr.Issues {
		log.Printf("%s: %s [%s]", issue.Key, issue.Fields.Summary, issue.Fields.Status.Name)
		s.repo.SaveIssue(issue)
	}

	return nil
}
