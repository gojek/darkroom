package gcs

import "github.com/gojektech/heimdall"

// Options represents the Google Cloud Storage storage options
type Options struct {
	// BucketName represents the name of the bucket
	BucketName string
	// CredentialsJSON holds the json data for credentials of a service account
	CredentialsJSON []byte
	// UseDefaultCredential toggle the usage of google application default credential to authenticate with cloud storage
	UseDefaultCredential bool
	// Client can be used to specify a heimdall.Client with hystrix like circuit breaker
	Client heimdall.Client
}
