package webfolder

import (
	"context"
	"fmt"
	"github.com/gojek/darkroom/pkg/storage"
	"github.com/gojektech/heimdall"
	"io/ioutil"
	"net/http"
)

// Storage holds the fields used by webfolder storage implementation
type Storage struct {
	baseURL string
	client  heimdall.Client
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the WebFolder storage backend.
func (s *Storage) Get(ctx context.Context, path string) storage.IResponse {
	res, err := s.client.Get(fmt.Sprintf("%s%s", s.baseURL, path), nil)
	if err != nil {
		if res != nil {
			return storage.NewResponse([]byte(nil), res.StatusCode, err)
		}
		return storage.NewResponse([]byte(nil), http.StatusUnprocessableEntity, err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	return storage.NewResponse([]byte(body), res.StatusCode, nil)
}

// NewStorage returns a new webfolder.Storage instance
func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
