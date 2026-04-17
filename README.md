# holidays.rest Go SDK

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/25967cd426ca4af5ac928747dbff939b)](https://app.codacy.com/gh/holidays-rest/sdk-go/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/25967cd426ca4af5ac928747dbff939b)](https://app.codacy.com/gh/holidays-rest/sdk-go/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

Official Go SDK for the [holidays.rest](https://www.holidays.rest) API.

## Requirements

- Go 1.21+
- No external dependencies — uses only the standard library

## Installation

```bash
go get github.com/holidays-rest/sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    holidays "github.com/holidays-rest/sdk-go"
)

func main() {
    client, err := holidays.NewClient("YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }

    h, err := client.Holidays(context.Background(), holidays.HolidaysParams{
        Country: "US",
        Year:    2024,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, holiday := range h {
        fmt.Printf("%s — %s\n", holiday.Date, holiday.Name["en"])
    }
}
```

Get an API key at [holidays.rest/dashboard](https://www.holidays.rest/dashboard).

---

## API

### `holidays.NewClient(apiKey, ...Option) (*Client, error)`

Creates a new client. Fails if `apiKey` is empty.

**Options:**

| Option                        | Description                                 |
|-------------------------------|---------------------------------------------|
| `WithBaseURL(url string)`     | Override base URL (useful for testing)      |
| `WithHTTPClient(*http.Client)`| Provide a custom HTTP client                |

```go
client, err := holidays.NewClient(
    "YOUR_API_KEY",
    holidays.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
```

---

### `client.Holidays(ctx, HolidaysParams) ([]Holiday, error)`

Fetch public holidays.

```go
type HolidaysParams struct {
    Country  string   // required — ISO 3166 alpha-2 (e.g. "US")
    Year     int      // required — e.g. 2024
    Month    int      // optional — 1–12
    Day      int      // optional — 1–31
    Type     []string // optional — "religious", "national", "local"
    Religion []int    // optional — religion codes 1–11
    Region   []string // optional — subdivision codes from Country()
    Lang     []string // optional — language codes from Languages()
    Response string   // optional — "json" (default) | "xml" | "yaml" | "csv"
}
```

Each returned `Holiday` has the following shape:

```go
type Holiday struct {
    CountryCode string            // ISO 3166 alpha-2, e.g. "DE"
    CountryName string            // e.g. "Germany"
    Date        string            // "YYYY-MM-DD"
    Name        map[string]string // language code → name, e.g. {"en": "New Year's Day"}
    IsNational  bool
    IsReligious bool
    IsLocal     bool
    IsEstimate  bool
    Day         struct {
        Actual   string // e.g. "Thursday"
        Observed string // e.g. "Thursday"
    }
    Religion string   // e.g. "Christianity", empty string if none
    Regions  []string // subdivision codes, e.g. ["BW", "BY"]
}
```

```go
// All US holidays in 2024
h, err := client.Holidays(ctx, holidays.HolidaysParams{
    Country: "US",
    Year:    2024,
})

// National holidays only
h, err := client.Holidays(ctx, holidays.HolidaysParams{
    Country: "DE",
    Year:    2024,
    Type:    []string{"national"},
})

// Multiple types
h, err := client.Holidays(ctx, holidays.HolidaysParams{
    Country: "TR",
    Year:    2024,
    Type:    []string{"national", "religious"},
})

// Filter by month and day
h, err := client.Holidays(ctx, holidays.HolidaysParams{
    Country: "GB",
    Year:    2024,
    Month:   12,
    Day:     25,
})

// Specific region
h, err := client.Holidays(ctx, holidays.HolidaysParams{
    Country: "US",
    Year:    2024,
    Region:  []string{"US-CA"},
})
```

---

### `client.Countries(ctx) ([]Country, error)`

List all supported countries.

```go
countries, err := client.Countries(ctx)
for _, c := range countries {
    fmt.Println(c.Alpha2, c.Name)
}
```

---

### `client.Country(ctx, countryCode) (*Country, error)`

Get country details including subdivision codes.

```go
us, err := client.Country(ctx, "US")
for _, s := range us.Subdivisions {
    fmt.Println(s.Code, s.Name)
}
```

---

### `client.Languages(ctx) ([]Language, error)`

List all supported language codes.

```go
langs, err := client.Languages(ctx)
```

---

## Error Handling

Non-2xx responses return `*APIError`:

```go
h, err := client.Holidays(ctx, holidays.HolidaysParams{Country: "US", Year: 2024})
if err != nil {
    var apiErr *holidays.APIError
    if errors.As(err, &apiErr) {
        fmt.Println(apiErr.Status)  // HTTP status code
        fmt.Println(apiErr.Message) // Error message
        fmt.Println(apiErr.Body)    // Raw response body
    }
    log.Fatal(err)
}
```

| Status | Meaning             |
|--------|---------------------|
| 400    | Bad request         |
| 401    | Invalid API key     |
| 404    | Not found           |
| 500    | Server error        |
| 503    | Service unavailable |

---

## License

MIT
