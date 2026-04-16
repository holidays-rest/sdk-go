package holidays

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCountries_DecodesResponse(t *testing.T) {
	fixture := []Country{
		{Name: "United States", Alpha2: "US", Subdivisions: []Subdivision{{Code: "US-NY", Name: "New York"}}},
		{Name: "Germany", Alpha2: "DE"},
	}
	body, _ := json.Marshal(fixture)

	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/countries" {
			t.Errorf("path = %q, want /countries", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	got, err := c.Countries(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Alpha2 != "US" {
		t.Errorf("Alpha2 = %q, want US", got[0].Alpha2)
	}
	if len(got[0].Subdivisions) != 1 || got[0].Subdivisions[0].Code != "US-NY" {
		t.Errorf("Subdivisions not decoded correctly")
	}
}

func TestCountry_EmptyCode(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	_, err := c.Country(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty countryCode")
	}
}

func TestCountry_DecodesResponse(t *testing.T) {
	fixture := Country{Name: "United States", Alpha2: "US"}
	body, _ := json.Marshal(fixture)

	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/country/US" {
			t.Errorf("path = %q, want /country/US", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	got, err := c.Country(context.Background(), "US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Alpha2 != "US" {
		t.Errorf("Alpha2 = %q, want US", got.Alpha2)
	}
}
