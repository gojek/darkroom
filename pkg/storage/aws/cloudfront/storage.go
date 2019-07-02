package cloudfront

import (
	"context"
	"fmt"
	"github.com/gojek/darkroom/pkg/storage"
	"github.com/gojektech/heimdall"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// Storage holds the fields used by cloudfront storage implementation
type Storage struct {
	cloudfrontHost string
	client         heimdall.Client
	secureProtocol bool
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the cloudfront storage backend.
func (s *Storage) Get(ctx context.Context, path string) storage.IResponse {
	res, err := s.client.Get(fmt.Sprintf("%s://%s%s", s.getProtocol(), s.cloudfrontHost, path), nil)
	if err != nil {
		if res != nil {
			return storage.NewResponse([]byte(nil), res.StatusCode, err)
		}
		return storage.NewResponse([]byte(nil), http.StatusUnprocessableEntity, err)
	}
	if res.StatusCode == http.StatusForbidden {
		return storage.NewResponse([]byte(nil), res.StatusCode, errors.New("forbidden"))
	}
	body, _ := ioutil.ReadAll(res.Body)
	return storage.NewResponse([]byte(body), res.StatusCode, nil)
}

func (s *Storage) getProtocol() string {
	if s.secureProtocol {
		return "https"
	}
	return "http"
}

// NewStorage returns a new cloudfront.Storage instance
func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
