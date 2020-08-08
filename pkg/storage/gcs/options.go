package gcs

// Options represents the Google Cloud Storage storage options
type Options struct {
	// BucketName represents the name of the bucket
	BucketName string
	// CredentialsJSON holds the json data for credentials of a service account
	CredentialsJSON []byte
}
