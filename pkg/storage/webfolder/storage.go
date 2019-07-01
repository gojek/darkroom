package webfolder

import (
	"context"
	"fmt"
	"github.com/gojektech/heimdall"
	"io/ioutil"
	"net/http"
	"github.com/gojek/darkroom/pkg/storage"
)

type Storage struct {
	baseURL string
	client  heimdall.Client
}

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

func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
