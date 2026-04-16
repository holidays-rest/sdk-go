package holidays

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// HolidayName maps language codes to localized holiday names (e.g. {"en": "New Year's Day"}).
type HolidayName map[string]string

// HolidayDay holds the actual and observed day-of-week for a holiday.
type HolidayDay struct {
	Actual   string `json:"actual"`
	Observed string `json:"observed"`
}

// Holiday represents a single public holiday returned by the API.
type Holiday struct {
	CountryCode string      `json:"country_code"`
	CountryName string      `json:"country_name"`
	Date        string      `json:"date"`
	Name        HolidayName `json:"name"`
	IsNational  bool        `json:"isNational"`
	IsReligious bool        `json:"isReligious"`
	IsLocal     bool        `json:"isLocal"`
	IsEstimate  bool        `json:"isEstimate"`
	Day         HolidayDay  `json:"day"`
	Religion    string      `json:"religion"`
	Regions     []string    `json:"regions"`
}

// HolidaysParams holds the parameters for the Holidays endpoint.
// Country and Year are required; all other fields are optional.
type HolidaysParams struct {
	// Required
	Country string // ISO 3166 alpha-2 code, e.g. "US"
	Year    int    // Four-digit year, e.g. 2024

	// Optional filters
	Month    int      // 1–12; zero value omits the parameter
	Day      int      // 1–31; zero value omits the parameter
	Type     []string // "religious", "national", "local"
	Religion []int    // Religion codes 1–11
	Region   []string // Region/subdivision codes from Country()
	Lang     []string // Language codes from Languages()
	Response string   // "json" (default) | "xml" | "yaml" | "csv"
}

// Holidays returns public holidays matching the given parameters.
func (c *Client) Holidays(ctx context.Context, p HolidaysParams) ([]Holiday, error) {
	if p.Country == "" {
		return nil, fmt.Errorf("holidays: Holidays: Country is required")
	}
	if p.Year == 0 {
		return nil, fmt.Errorf("holidays: Holidays: Year is required")
	}

	q := url.Values{}
	q.Set("country", p.Country)
	q.Set("year", strconv.Itoa(p.Year))

	if p.Month != 0 {
		q.Set("month", strconv.Itoa(p.Month))
	}
	if p.Day != 0 {
		q.Set("day", strconv.Itoa(p.Day))
	}
	if len(p.Type) > 0 {
		q.Set("type", strings.Join(p.Type, ","))
	}
	if len(p.Religion) > 0 {
		codes := make([]string, len(p.Religion))
		for i, r := range p.Religion {
			codes[i] = strconv.Itoa(r)
		}
		q.Set("religion", strings.Join(codes, ","))
	}
	if len(p.Region) > 0 {
		q.Set("region", strings.Join(p.Region, ","))
	}
	if len(p.Lang) > 0 {
		q.Set("lang", strings.Join(p.Lang, ","))
	}
	if p.Response != "" {
		q.Set("response", p.Response)
	}

	var result []Holiday
	if err := c.get(ctx, "/holidays", q, &result); err != nil {
		return nil, err
	}

	return result, nil
}
