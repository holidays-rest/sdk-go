package holidays

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := NewClient("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return srv, c
}

func TestNewClient_EmptyAPIKey(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Fatal("expected error for empty apiKey, got nil")
	}
}

func TestNewClient_ValidAPIKey(t *testing.T) {
	c, err := NewClient("my-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.apiKey != "my-key" {
		t.Fatalf("apiKey = %q, want %q", c.apiKey, "my-key")
	}
	if c.baseURL != defaultBaseURL {
		t.Fatalf("baseURL = %q, want %q", c.baseURL, defaultBaseURL)
	}
}

func TestWithBaseURL_TrailingSlashStripped(t *testing.T) {
	c, _ := NewClient("k", WithBaseURL("https://example.com/v1/"))
	if c.baseURL != "https://example.com/v1" {
		t.Fatalf("baseURL = %q, trailing slash not stripped", c.baseURL)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c, _ := NewClient("k", WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Fatal("custom http client not set")
	}
}

func TestGet_SetsAuthorizationHeader(t *testing.T) {
	var gotAuth string
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})

	c.Languages(context.Background()) //nolint:errcheck

	if gotAuth != "Bearer test-key" {
		t.Fatalf("Authorization = %q, want %q", gotAuth, "Bearer test-key")
	}
}

func TestGet_APIError_WithJSONMessage(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"invalid api key"}`))
	})

	_, err := c.Languages(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.Status != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", apiErr.Status, http.StatusUnauthorized)
	}
	if apiErr.Message != "invalid api key" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "invalid api key")
	}
}

func TestGet_APIError_NonJSONBody(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
	})

	_, err := c.Languages(context.Background())
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Status != http.StatusInternalServerError {
		t.Errorf("Status = %d, want 500", apiErr.Status)
	}
}

func TestGet_InvalidJSON(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not-json`))
	})

	_, err := c.Languages(context.Background())
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}
