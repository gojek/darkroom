package gcs

import (
	"context"
	"net/http"

	"cloud.google.com/go/storage"

	"google.golang.org/api/option"

	"github.com/gojektech/heimdall"
	gcloud "google.golang.org/api/transport/http"
)

const userAgent = "gcloud-golang-storage/20151204"

func newHeimdallHTTPClient(ctx context.Context, hc heimdall.Client, credentialsJson []byte) (*http.Client, error) {
	t, err := newTransport(ctx, hc, credentialsJson)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: t,
	}, nil
}

func newTransport(ctx context.Context, hc heimdall.Client, credentialsJson []byte) (http.RoundTripper, error) {
	o := option.WithoutAuthentication()
	if len(credentialsJson) > 0 {
		o = option.WithCredentialsJSON(credentialsJson)
	}
	return gcloud.NewTransport(ctx,
		&hystrixTransport{client: hc},
		option.WithUserAgent(userAgent),
		option.WithScopes(storage.ScopeReadOnly),
		o)
}

type hystrixTransport struct {
	client heimdall.Client
}

func (h hystrixTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return h.client.Do(request)
}
