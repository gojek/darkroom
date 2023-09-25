package gcs

import (
	"context"
	"golang.org/x/oauth2/google"
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

func newTransport(ctx context.Context, opts *Options) (http.RoundTripper, error) {
	o := option.WithoutAuthentication()
	if opts.UseDefaultCredential {
		credential, err := google.FindDefaultCredentials(ctx)
		if err != nil {
			return nil, err
		}
		o = option.WithCredentials(credential)
	} else if len(opts.CredentialsJSON) > 0 {
		o = option.WithCredentialsJSON(opts.CredentialsJSON)
	}
	return gcloud.NewTransport(ctx,
		&hystrixTransport{client: opts.Client},
		option.WithUserAgent(userAgent),
		option.WithScopes(storage.ScopeReadOnly),
		o,
	)

}

type hystrixTransport struct {
	client heimdall.Client
}

func (h hystrixTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return h.client.Do(request)
}
