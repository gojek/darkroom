package gcs

import (
	"net/http"

	"github.com/gojektech/heimdall"
)

func newHeimdallHTTPClient(hc heimdall.Client) *http.Client {
	return &http.Client{
		Transport: newTransport(hc),
	}
}

func newTransport(hc heimdall.Client) http.RoundTripper {
	return hystrixTransport{client: hc}
}

type hystrixTransport struct {
	client heimdall.Client
}

func (h hystrixTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return h.client.Do(request)
}
