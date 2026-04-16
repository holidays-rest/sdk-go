package holidays

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func TestHolidays_MissingCountry(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	_, err := c.Holidays(context.Background(), HolidaysParams{Year: 2024})
	if err == nil {
		t.Fatal("expected error for missing Country")
	}
}

func TestHolidays_MissingYear(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	_, err := c.Holidays(context.Background(), HolidaysParams{Country: "US"})
	if err == nil {
		t.Fatal("expected error for missing Year")
	}
}

func TestHolidays_RequiredQueryParams(t *testing.T) {
	var gotQuery url.Values
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})

	_, err := c.Holidays(context.Background(), HolidaysParams{Country: "US", Year: 2024})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery.Get("country") != "US" {
		t.Errorf("country = %q, want US", gotQuery.Get("country"))
	}
	if gotQuery.Get("year") != "2024" {
		t.Errorf("year = %q, want 2024", gotQuery.Get("year"))
	}
}

func TestHolidays_OptionalQueryParams(t *testing.T) {
	var gotQuery url.Values
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})

	_, err := c.Holidays(context.Background(), HolidaysParams{
		Country:  "US",
		Year:     2024,
		Month:    12,
		Day:      25,
		Type:     []string{"national", "religious"},
		Religion: []int{1, 2},
		Region:   []string{"US-NY"},
		Lang:     []string{"en"},
		Response: "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := map[string]string{
		"month":    "12",
		"day":      "25",
		"type":     "national,religious",
		"religion": "1,2",
		"region":   "US-NY",
		"lang":     "en",
		"response": "json",
	}
	for param, want := range cases {
		if got := gotQuery.Get(param); got != want {
			t.Errorf("param %q = %q, want %q", param, got, want)
		}
	}
}

func TestHolidays_OptionalParamsOmittedWhenZero(t *testing.T) {
	var gotQuery url.Values
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})

	c.Holidays(context.Background(), HolidaysParams{Country: "US", Year: 2024}) //nolint:errcheck

	for _, param := range []string{"month", "day", "type", "religion", "region", "lang", "response"} {
		if gotQuery.Has(param) {
			t.Errorf("param %q should be absent when zero/empty", param)
		}
	}
}

func TestHolidays_DecodesResponse(t *testing.T) {
	fixture := []Holiday{
		{
			CountryCode: "US",
			CountryName: "United States",
			Date:        "2024-12-25",
			Name:        HolidayName{"en": "Christmas Day"},
			IsNational:  true,
			Day:         HolidayDay{Actual: "Wednesday", Observed: "Wednesday"},
		},
	}
	body, _ := json.Marshal(fixture)

	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	got, err := c.Holidays(context.Background(), HolidaysParams{Country: "US", Year: 2024})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].CountryCode != "US" {
		t.Errorf("CountryCode = %q, want US", got[0].CountryCode)
	}
	if got[0].Name["en"] != "Christmas Day" {
		t.Errorf("Name[en] = %q, want Christmas Day", got[0].Name["en"])
	}
}
