package jira

type JiraClient struct{}

func New() (*JiraClient, error) {
	return &JiraClient{}, nil
}
