package main

import (
	"errors"
	"strings"
	"testing"
)

func TestParseError(t *testing.T) {
	body := []byte(`{"error":{"message":"Invalid OAuth access token.","type":"OAuthException","code":190,"error_subcode":460,"error_data":{"a":1},"fbtrace_id":"abc"}}`)
	err := parseError(400, body)

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("%T is not *Error", err)
	}
	if apiErr.Code != 190 || apiErr.Subcode != 460 || apiErr.TraceID != "abc" || apiErr.StatusCode != 400 {
		t.Errorf("unexpected fields: %+v", apiErr)
	}
	if want := "Error occured: Invalid OAuth access token."; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestParseErrorNonJSON(t *testing.T) {
	err := parseError(502, []byte("<html>Bad Gateway</html>"))
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 502 {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "HTTP 502") || !strings.Contains(err.Error(), "Bad Gateway") {
		t.Errorf("message should contain status and body: %q", err.Error())
	}
}
