package gcs

import (
	gs "cloud.google.com/go/storage"
	"context"
	"errors"
	"github.com/gojek/darkroom/pkg/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
)

// Storage holds the fields used by Google Cloud Storage implementation
type Storage struct {
	bucketHandle BucketHandle
}

// NewStorage returns a new gcs.Storage instance
func NewStorage(opts Options) (*Storage, error) {
	c, err := gs.NewClient(context.TODO(), clientOptions(opts.CredentialsJSON)...)
	if err != nil {
		return nil, err
	}
	return &Storage{bucketHandle{c.Bucket(opts.BucketName)}}, nil
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the S3 storage backend.
func (s *Storage) Get(ctx context.Context, path string) storage.IResponse {
	r, err := s.bucketHandle.Object(path).NewReader(ctx)
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		return storage.NewResponse(nil, apiErr.Code, apiErr)
	}
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return storage.NewResponse(nil, http.StatusUnprocessableEntity, err)
	}
	return storage.NewResponse(d, http.StatusOK, nil)
}

func clientOptions(credentialsJSON []byte) []option.ClientOption {
	if len(credentialsJSON) != 0 {
		return []option.ClientOption{option.WithCredentialsJSON(credentialsJSON)}
	}
	return []option.ClientOption{option.WithoutAuthentication()}
}
