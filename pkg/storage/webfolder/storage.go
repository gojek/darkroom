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
func (s *Storage) Get(ctx context.Context, path string, opt *storage.GetRequestOptions) storage.IResponse {
	var h http.Header
	if opt != nil && opt.Range != "" {
		h = http.Header{}
		h.Add(storage.HeaderRange, opt.Range)
	}

	res, err := s.client.Get(fmt.Sprintf("%s%s", s.baseURL, path), h)
	if err != nil {
		if res != nil {
			return storage.NewResponse([]byte(nil), res.StatusCode, err, nil)
		}
		return storage.NewResponse([]byte(nil), http.StatusUnprocessableEntity, err, nil)
	}
	body, _ := ioutil.ReadAll(res.Body)
	return storage.NewResponse(body, res.StatusCode, nil, s.newMetadata(&res.Header))
}

func (s *Storage) newMetadata(header *http.Header) *storage.ResponseMetadata {
	return &storage.ResponseMetadata{
		AcceptRanges:  header.Get(storage.HeaderAcceptRanges),
		ContentLength: header.Get(storage.HeaderContentLength),
		ContentRange:  header.Get(storage.HeaderContentRange),
		ContentType:   header.Get(storage.HeaderContentType),
		ETag:          header.Get(storage.HeaderETag),
		LastModified:  header.Get(storage.HeaderLastModified),
	}
}

// NewStorage returns a new webfolder.Storage instance
func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
