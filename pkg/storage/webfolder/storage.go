package webfolder

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gojek/darkroom/pkg/storage"
	"github.com/gojektech/heimdall"
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
	if resErr, ok := s.hasError(res, err); ok {
		return resErr
	}

	body, _ := ioutil.ReadAll(res.Body)
	return storage.NewResponse(body, res.StatusCode, nil)
}

// GetPartially takes in the Context, path and opt (which ignored) as an argument and returns an IResponse interface implementation.
// This method is an alias of `Get` method
func (s *Storage) GetPartially(ctx context.Context, path string, _ *storage.GetPartiallyRequestOptions) storage.IResponse {
	return s.Get(ctx, path)
}

func (s *Storage) hasError(res *http.Response, err error) (storage.IResponse, bool) {
	if err != nil {
		if res != nil {
			return storage.NewResponse([]byte(nil), res.StatusCode, err), true
		}
		return storage.NewResponse([]byte(nil), http.StatusUnprocessableEntity, err), true
	}
	return nil, false
}

// NewStorage returns a new webfolder.Storage instance
func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
