package holidays

import "fmt"

// APIError is returned when the server responds with a non-2xx status code.
type APIError struct {
	Status  int
	Message string
	Body    []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("holidays.rest API error %d: %s", e.Status, e.Message)
}
