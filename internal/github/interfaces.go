package github


import "net/http"

// HTTPClient a service that can do http calls
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
