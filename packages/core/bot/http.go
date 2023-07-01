package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var errHTTP = errors.New("http error")

const (
	bearer        = "Bearer"
	authorization = "authorization"

	accept      = "Accept"
	contentType = "Content-Type"

	jsonContentType = "application/json"
)

func NewJSONRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w. %w", errHTTP, err)
	}

	req.Header.Set(accept, jsonContentType)

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set(contentType, jsonContentType)
	}

	return req, nil
}
