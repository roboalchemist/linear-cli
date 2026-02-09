package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseURL = "https://api.linear.app/graphql"
)

type Client struct {
	httpClient    *http.Client
	authHeader    string
	baseURL       string
	LastRateLimit *RateLimit // Updated after each request
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLErrorLocation `json:"locations,omitempty"`
	Path      []interface{}          `json:"path,omitempty"`
}

type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// NewClient creates a new Linear API client
func NewClient(authHeader string) *Client {
	return NewClientWithURL(BaseURL, authHeader)
}

// NewClientWithURL creates a new Linear API client with custom URL
func NewClientWithURL(baseURL, authHeader string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		authHeader: authHeader,
		baseURL:    baseURL,
	}
}

// Execute performs a GraphQL request
func (c *Client) Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("User-Agent", "linear-cli/0.1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Capture rate limit headers
	c.LastRateLimit = parseRateLimit(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("GraphQL errors: %v", gqlResp.Errors)
	}

	if result != nil {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}

// ExecuteRaw performs a GraphQL request and returns the raw JSON response data
func (c *Client) ExecuteRaw(ctx context.Context, query string, variables map[string]interface{}) (json.RawMessage, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("User-Agent", "linear-cli/0.1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Capture rate limit headers
	c.LastRateLimit = parseRateLimit(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", gqlResp.Errors)
	}

	return gqlResp.Data, nil
}

// RateLimit holds rate limit info parsed from Linear API response headers
type RateLimit struct {
	// Request limits
	RequestLimit     int       `json:"requestLimit"`
	RequestRemaining int       `json:"requestRemaining"`
	RequestReset     time.Time `json:"requestReset"`
	// Complexity limits
	Complexity          int       `json:"complexity"`
	ComplexityLimit     int       `json:"complexityLimit"`
	ComplexityRemaining int       `json:"complexityRemaining"`
	ComplexityReset     time.Time `json:"complexityReset"`
}

// parseRateLimit extracts rate limit info from HTTP response headers
func parseRateLimit(resp *http.Response) *RateLimit {
	rl := &RateLimit{}
	rl.RequestLimit = headerInt(resp, "X-RateLimit-Requests-Limit")
	rl.RequestRemaining = headerInt(resp, "X-RateLimit-Requests-Remaining")
	rl.RequestReset = headerTime(resp, "X-RateLimit-Requests-Reset")
	rl.Complexity = headerInt(resp, "X-Complexity")
	rl.ComplexityLimit = headerInt(resp, "X-RateLimit-Complexity-Limit")
	rl.ComplexityRemaining = headerInt(resp, "X-RateLimit-Complexity-Remaining")
	rl.ComplexityReset = headerTime(resp, "X-RateLimit-Complexity-Reset")
	return rl
}

func headerInt(resp *http.Response, key string) int {
	v := resp.Header.Get(key)
	if v == "" {
		return 0
	}
	n, _ := strconv.Atoi(v)
	return n
}

func headerTime(resp *http.Response, key string) time.Time {
	v := resp.Header.Get(key)
	if v == "" {
		return time.Time{}
	}
	ms, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.UnixMilli(ms)
}

// GetRateLimit makes a lightweight request and returns the current rate limit status
func (c *Client) GetRateLimit(ctx context.Context) (*RateLimit, error) {
	// Use a minimal query to get rate limit headers
	err := c.Execute(ctx, `query { viewer { id } }`, nil, nil)
	if err != nil {
		return nil, err
	}
	if c.LastRateLimit == nil {
		return nil, fmt.Errorf("no rate limit info available")
	}
	return c.LastRateLimit, nil
}
