package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JiraClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func New(baseURL, token string) *JiraClient {
	return &JiraClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *JiraClient) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return c.client.Do(req.WithContext(ctx))
}

type Issue struct {
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary string `json:"summary"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}

type SearchResponse struct {
	Issues []Issue `json:"issues"`
	Total  int     `json:"total"`
}

func (c *JiraClient) GetIssue(key string) (*Issue, error) {
	resp, err := c.do(context.Background(), "GET", fmt.Sprintf("/rest/api/2/issue/%s", key), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Jira API error: %d, body=%s", resp.StatusCode, string(body))
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *JiraClient) SearchIssues(jql string) (*SearchResponse, error) {
	encodedJQL := url.QueryEscape(jql)
	fields := "summary,status"
	urlStr := fmt.Sprintf("/rest/api/2/search?jql=%s&fields=%s", encodedJQL, fields)

	resp, err := c.do(context.Background(), "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed: %d, %s", resp.StatusCode, string(body))
	}

	var sr SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return &sr, nil
}
