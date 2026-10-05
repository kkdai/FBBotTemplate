package main

import (
	"bytes"
	"io"
	"net/http"
)

type staticTransport struct {
	statusCode int
	body       []byte
}

func (t staticTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: t.statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(t.body)),
		Request:    req,
	}, nil
}

func setClient(statusCode int, body []byte) {
	httpClient = &http.Client{
		Transport: staticTransport{
			statusCode: statusCode,
			body:       body,
		},
	}
}
