package messenger

import (
	"encoding/json"
	"fmt"
	"strings"
)

type rawError struct {
	Error Error `json:"error"`
}

// Error is an error returned by the Graph API. Use errors.As to inspect Code and TraceID.
type Error struct {
	Message   string          `json:"message"`
	Type      string          `json:"type"`
	Code      int             `json:"code"`
	Subcode   int             `json:"error_subcode,omitempty"`
	ErrorData json.RawMessage `json:"error_data,omitempty"`
	TraceID   string          `json:"fbtrace_id"`
	// StatusCode is the HTTP status of the response.
	StatusCode int `json:"-"`
}

func (e *Error) Error() string {
	return "Error occured: " + e.Message
}

// parseError converts a non-200 response body into an *Error.
// If the body is not a Graph API error, the HTTP status and raw body are used as the message.
func parseError(statusCode int, body []byte) error {
	raw := new(rawError)
	if err := json.Unmarshal(body, raw); err != nil || raw.Error.Message == "" {
		raw.Error = Error{Message: fmt.Sprintf("HTTP %d: %s", statusCode, strings.TrimSpace(string(body)))}
	}
	raw.Error.StatusCode = statusCode
	return &raw.Error
}
