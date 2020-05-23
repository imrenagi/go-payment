// +build !go1.13

package xendit

import (
	"context"
	"io"
	"net/http"
)

func newHTTPRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}

	return req.WithContext(ctx), nil
}
