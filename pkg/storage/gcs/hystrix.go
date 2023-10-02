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

func newHeimdallHTTPClient(ctx context.Context, opts *Options) (*http.Client, error) {
	t, err := newTransport(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: t,
	}, nil
}

func getCredentialsOption(opts *Options) option.ClientOption {
	if opts.Credentials != nil {
		return option.WithCredentials(opts.Credentials)
	}
	if len(opts.CredentialsJSON) > 0 {
		return option.WithCredentialsJSON(opts.CredentialsJSON)
	}
	return option.WithoutAuthentication()
}

func newTransport(ctx context.Context, opts *Options) (http.RoundTripper, error) {
	return gcloud.NewTransport(ctx,
		&hystrixTransport{client: opts.Client},
		option.WithUserAgent(userAgent),
		option.WithScopes(storage.ScopeReadOnly),
		getCredentialsOption(opts),
	)
}

type hystrixTransport struct {
	client heimdall.Client
}

func (h hystrixTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return h.client.Do(request)
}
