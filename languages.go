package holidays

import (
	"context"
	"net/url"
)

// Language represents a supported language.
type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Languages returns all supported language codes.
func (c *Client) Languages(ctx context.Context) ([]Language, error) {
	var result []Language
	if err := c.get(ctx, "/languages", url.Values{}, &result); err != nil {
		return nil, err
	}
	return result, nil
}
