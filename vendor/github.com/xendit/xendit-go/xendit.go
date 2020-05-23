// Package xendit provides the binding for Xendit APIs.
package xendit

import (
	"net/http"
	"sync"
)

// Opt is the default Option for the API call without API client
var Opt Option = Option{
	XenditURL: "https://api.xendit.co",
}

var apiRequesterWrapper APIRequesterWrapper = APIRequesterWrapper{}

var httpClient *http.Client = &http.Client{}

// Option is the wrap of the parameters needed for the API call
type Option struct {
	SecretKey string // customer's secret API key
	XenditURL string // should there be a need to override API base URL
}

// APIRequesterWrapper is the APIRequester with locker for setting the APIRequester
type APIRequesterWrapper struct {
	apiRequester APIRequester
	mu           sync.RWMutex
}

// GetAPIRequester returns the xendit APIRequester.
// If it is already created, it will return the created one.
// Else, it will create a default implementation.
func GetAPIRequester() APIRequester {
	if apiRequesterWrapper.apiRequester != nil {
		return apiRequesterWrapper.apiRequester
	}

	apiRequesterWrapper.apiRequester = &APIRequesterImplementation{
		HTTPClient: httpClient,
	}

	return apiRequesterWrapper.apiRequester
}

// SetAPIRequester sets the APIRequester for API call
func SetAPIRequester(apiRequester APIRequester) {
	apiRequesterWrapper.mu.Lock()
	defer apiRequesterWrapper.mu.Unlock()

	apiRequesterWrapper.apiRequester = apiRequester
}

// SetHTTPClient sets the httpClient for API call
func SetHTTPClient(newHTTPClient *http.Client) {
	httpClient = newHTTPClient
}
