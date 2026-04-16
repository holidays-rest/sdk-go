package holidays

import "testing"

func TestAPIError_Error(t *testing.T) {
	e := &APIError{Status: 404, Message: "not found"}
	want := "holidays.rest API error 404: not found"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
