package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	defaultTimeout    = 30 * time.Second
	rateLimitInterval = 100 * time.Millisecond
	defaultRetryAfter = 30 * time.Second
)

type Client struct {
	baseURL string
	token   string
	client  *http.Client
	limiter *rate.Limiter
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

type Error struct {
	StatusCode    int               `json:"-"`
	Body          []byte            `json:"-"`
	Message       string            `json:"-"`
	ErrorMessages []string          `json:"errorMessages,omitempty"`
	Errors        map[string]string `json:"errors,omitempty"`
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Jira API: %d, message: %s", e.StatusCode, e.Message)
	}

	if len(e.ErrorMessages) > 0 {
		return fmt.Sprintf("Jira API: %d, errors: %v", e.StatusCode, e.ErrorMessages)
	}

	if len(e.Errors) > 0 {
		return fmt.Sprintf("Jira API: %d, field errors: %v", e.StatusCode, e.Errors)
	}

	return fmt.Sprintf("Jira API: %d, body: %s", e.StatusCode, string(e.Body))
}

func (e *Error) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

//nolint:unused
func getRetryAfter(resp *http.Response) time.Duration {
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		if sec, err := strconv.Atoi(retryAfter); err == nil {
			return time.Duration(sec) * time.Second
		}
	}

	return defaultRetryAfter
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		limiter: rate.NewLimiter(rate.Every(rateLimitInterval), 1),
	}
}

func (c *Client) GetIssue(ctx context.Context, key string) (*Issue, error) {
	resp, err := c.do(ctx, http.MethodGet, "/rest/api/2/issue/"+key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue %q: %w", key, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("decode Jira issue response: %w", err)
	}

	return &issue, nil
}

func (c *Client) SearchIssues(ctx context.Context, jql string) (*SearchResponse, error) {
	encodedJQL := url.QueryEscape(jql)
	fields := "summary,status"
	urlStr := fmt.Sprintf("/rest/api/2/search?jql=%s&fields=%s", encodedJQL, fields)

	resp, err := c.do(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET request to %q: %w", urlStr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("decode Jira search response: %w", err)
	}

	return &searchResp, nil
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response body: %w", err)
	}

	errAPI := &Error{
		StatusCode: resp.StatusCode,
		Body:       body,
	}

	if err := json.Unmarshal(body, errAPI); err != nil {
		errAPI.Message = string(body)
	}

	return errAPI
}

func (c *Client) do(
	ctx context.Context,
	method, path string,
	body io.Reader,
) (*http.Response, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter wait failed: %w", err)
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("execute HTTP request: %w", err)
	}

	return resp, nil
}
