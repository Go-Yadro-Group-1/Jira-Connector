package jira

//nolint:revive
type JiraClient struct{}

func New() (*JiraClient, error) {
	return &JiraClient{}, nil
}
