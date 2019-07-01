package cloudfront

import (
	"context"
	"fmt"
	"github.com/gojektech/heimdall"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"github.com/gojek/darkroom/pkg/storage"
)

type Storage struct {
	cloudfrontHost string
	client         heimdall.Client
	secureProtocol bool
}

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

func NewStorage(opts ...Option) *Storage {
	s := Storage{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
