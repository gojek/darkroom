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
	cloudfrontHost string // can end with trailing slash or not (example: "localhost:8000", "localhost:8000/")
	client         heimdall.Client
	secureProtocol bool
}

// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
// This method figures out how to get the data from the cloudfront storage backend.
func (s *Storage) Get(ctx context.Context, path string) storage.IResponse {
	res, err := s.client.Get(s.getURL(path), nil)
	if errRes, ok := s.hasError(res, err); ok {
		return errRes
	}
	body, _ := ioutil.ReadAll(res.Body)
	return storage.NewResponse(body, res.StatusCode, nil)
}

// GetPartialObject takes in the Context, path and opt as an argument and returns an IResponse interface implementation.
// This method figures out how to get partial data from the cloudfront storage backend.
func (s *Storage) GetPartialObject(ctx context.Context, path string, opt *storage.GetPartialObjectRequestOptions) storage.IResponse {
	var h http.Header
	if opt != nil && opt.Range != "" {
		h = http.Header{}
		h.Add(storage.HeaderRange, opt.Range)
	}

	res, err := s.client.Get(s.getURL(path), h)
	if errRes, ok := s.hasError(res, err); ok {
		return errRes
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

func (s *Storage) getURL(path string) string {
	host := s.cloudfrontHost
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	if path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s://%s%s", s.getProtocol(), host, path)
}

func (s *Storage) getProtocol() string {
	if s.secureProtocol {
		return "https"
	}
	return "http"
}

func (s *Storage) hasError(res *http.Response, err error) (storage.IResponse, bool) {
	if err != nil {
		if res != nil {
			return storage.NewResponse([]byte(nil), res.StatusCode, err), true
		}
		return storage.NewResponse([]byte(nil), http.StatusUnprocessableEntity, err), true
	}
	if res.StatusCode == http.StatusForbidden {
		return storage.NewResponse([]byte(nil), res.StatusCode, errors.New("forbidden")), true
	}
	return nil, false
}

// NewStorage returns a new cloudfront.Storage instance
func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
