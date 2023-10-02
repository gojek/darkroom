package gcs

import (
	"github.com/gojektech/heimdall"
	"golang.org/x/oauth2/google"
)

// Options represents the Google Cloud Storage storage options
type Options struct {
	// BucketName represents the name of the bucket
	BucketName string
	// CredentialsJSON holds the json data for credentials of a service account
	CredentialsJSON []byte
	// Credentials represents google credentials, including Application Default Credentials
	Credentials *google.Credentials
	// Client can be used to specify a heimdall.Client with hystrix like circuit breaker
	Client heimdall.Client
}
