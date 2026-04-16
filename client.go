// Package holidays provides a client for the holidays.rest API.
// Documentation: https://docs.holidays.rest
package holidays

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

const defaultBaseURL = "https://api.holidays.rest/v1"

// Client is the holidays.rest API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL overrides the API base URL. Useful for testing.
func WithBaseURL(u string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(u, "/")
	}
}

// WithHTTPClient replaces the default HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// NewClient creates a new Client. apiKey is required; obtain one at
// https://www.holidays.rest/dashboard.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("holidays: apiKey must not be empty")
	}

	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	for _, o := range opts {
		o(c)
	}

	return c, nil
}

// get executes a GET request, decodes JSON into dst, and returns any error.
func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("holidays: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("holidays: http: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("holidays: read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{Status: resp.StatusCode, Body: body}

		// Try to extract message from JSON body.
		var errBody struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &errBody) == nil && errBody.Message != "" {
			apiErr.Message = errBody.Message
		} else {
			apiErr.Message = resp.Status
		}

		return apiErr
	}

	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("holidays: decode response: %w", err)
	}

	return nil
}
