package gql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	URL        string
	HTTPClient *http.Client
	Headers    map[string]string
}

// NewClient returns a new GraphQL client with default settings.
func NewClient(url string) *Client {
	return &Client{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Headers: make(map[string]string),
	}
}

// SetHeader sets a custom header for all requests.
func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

// GraphQLRequest is the payload for a GraphQL request.
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse is the response from a GraphQL request.
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error.
type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []GraphQLLocation      `json:"locations,omitempty"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// GraphQLLocation indicates the location of a GraphQL error.
type GraphQLLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Error implements the error interface for GraphQLError.
func (e GraphQLError) Error() string {
	return e.Message
}

// Query performs a POST request with the given GraphQL query string and decodes the result into out.
func Query(ctx context.Context, url string, query string, out any) error {
	client := NewClient(url)
	return client.Query(ctx, query, nil, out)
}

// Query executes a GraphQL query with optional variables.
func (c *Client) Query(ctx context.Context, query string, variables map[string]interface{}, out any) error {
	return c.execute(ctx, query, variables, out)
}

// QueryWithVariables executes a GraphQL query with variables.
func (c *Client) QueryWithVariables(ctx context.Context, query string, variables map[string]interface{}, out any) error {
	return c.execute(ctx, query, variables, out)
}

func (c *Client) execute(ctx context.Context, query string, variables map[string]interface{}, out any) error {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("graphql request failed with status %s", resp.Status)
	}

	var gqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("graphql errors: %v", gqlResp.Errors)
	}

	if out != nil && gqlResp.Data != nil {
		if err := json.Unmarshal(gqlResp.Data, out); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}
