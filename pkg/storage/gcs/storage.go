package gcs

import (
	gs "cloud.google.com/go/storage"
	"context"
	"errors"
	"github.com/gojek/darkroom/pkg/storage"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"net/http"
)

// Storage holds the fields used by Google Cloud Storage implementation
type Storage struct {
	bucketHandle BucketHandle
}

// NewStorage returns a new gcs.Storage instance
func NewStorage(opts Options) *Storage {
	// TODO: Handle error
	c, _ := gs.NewClient(context.TODO())
	return &Storage{bucketHandle{c.Bucket(opts.BucketName)}}
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the S3 storage backend.
func (s *Storage) Get(ctx context.Context, path string) storage.IResponse {
	r, err := s.bucketHandle.Object(path).NewReader(ctx)
	apiErr := &googleapi.Error{}
	if errors.As(err, &apiErr) {
		return storage.NewResponse(nil, apiErr.Code, apiErr)
	}
	d, _ := ioutil.ReadAll(r)
	// TODO: Handle error
	return storage.NewResponse(d, http.StatusOK, nil)
}
