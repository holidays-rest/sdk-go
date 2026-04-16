package holidays

import (
	"context"
	"fmt"
	"net/url"
)

// Subdivision represents a region or administrative subdivision of a country.
type Subdivision struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Country represents a supported country.
type Country struct {
	Name         string        `json:"name"`
	Alpha2       string        `json:"alpha2"`
	Subdivisions []Subdivision `json:"subdivisions,omitempty"`
}

// Countries returns all supported countries.
func (c *Client) Countries(ctx context.Context) ([]Country, error) {
	var result []Country
	if err := c.get(ctx, "/countries", url.Values{}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Country returns details for a single country identified by its ISO 3166
// alpha-2 code (e.g. "US"). The response includes available subdivision codes
// that can be used as region filters in Holidays().
func (c *Client) Country(ctx context.Context, countryCode string) (*Country, error) {
	if countryCode == "" {
		return nil, fmt.Errorf("holidays: Country: countryCode must not be empty")
	}

	var result Country
	if err := c.get(ctx, "/country/"+countryCode, url.Values{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
