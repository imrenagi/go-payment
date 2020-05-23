// +build go1.13

package xendit

import (
	"context"
	"io"
	"net/http"
)

func newHTTPRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(
		ctx,
		method,
		url,
		body,
	)
}
