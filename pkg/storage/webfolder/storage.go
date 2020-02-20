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
	if resErr, ok := s.hasError(res, err); ok {
		return resErr
	}

	body, _ := ioutil.ReadAll(res.Body)
	return storage.NewResponse(body, res.StatusCode, nil)
}

// GetPartialObject takes in the Context, path and opt as an argument and returns an IResponse interface implementation.
// This method figures out how to get partial data from the WebFolder storage backend.
func (s *Storage) GetPartialObject(ctx context.Context, path string, opt *storage.GetPartialObjectRequestOptions) storage.IResponse {
	var h http.Header
	if opt != nil && opt.Range != "" {
		h = http.Header{}
		h.Add(storage.HeaderRange, opt.Range)
	}

	res, err := s.client.Get(fmt.Sprintf("%s%s", s.baseURL, path), h)
	if resErr, ok := s.hasError(res, err); ok {
		return resErr
	}

	body, _ := ioutil.ReadAll(res.Body)
	return storage.
		NewResponse(body, res.StatusCode, nil).
		WithMetadata(&storage.ResponseMetadata{
			AcceptRanges:  res.Header.Get(storage.HeaderAcceptRanges),
			ContentLength: res.Header.Get(storage.HeaderContentLength),
			ContentRange:  res.Header.Get(storage.HeaderContentRange),
			ContentType:   res.Header.Get(storage.HeaderContentType),
			ETag:          res.Header.Get(storage.HeaderETag),
			LastModified:  res.Header.Get(storage.HeaderLastModified),
		})
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
