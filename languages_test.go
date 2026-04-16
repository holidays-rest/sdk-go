package holidays

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestLanguages_DecodesResponse(t *testing.T) {
	fixture := []Language{
		{Code: "en", Name: "English"},
		{Code: "de", Name: "German"},
	}
	body, _ := json.Marshal(fixture)

	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/languages" {
			t.Errorf("path = %q, want /languages", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	got, err := c.Languages(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Code != "en" {
		t.Errorf("Code = %q, want en", got[0].Code)
	}
}

func TestLanguages_EmptyList(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})

	got, err := c.Languages(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("len = %d, want 0", len(got))
	}
}
