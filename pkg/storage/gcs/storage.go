package gcs

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	gs "cloud.google.com/go/storage"
	"github.com/gojek/darkroom/pkg/storage"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

var (
	rangeRegex      = regexp.MustCompile(`bytes=(\d+)-(\d+)`)
	ErrInvalidRange = errors.New("invalid range")
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
	path = strings.TrimPrefix(path, "/")
	r, err := s.bucketHandle.Object(path).NewReader(ctx)
	if errors.Is(err, gs.ErrObjectNotExist) {
		return storage.NewResponse(nil, http.StatusNotFound, err)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		return storage.NewResponse(nil, apiErr.Code, apiErr)
	}
	defer r.Close()
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return storage.NewResponse(nil, http.StatusUnprocessableEntity, err)
	}
	return storage.NewResponse(d, http.StatusOK, nil)
}

// GetPartially takes in the Context, path and opt as an argument and returns an IResponse interface implementation.
// This method figures out how to get partial data from the S3 storage backend.
func (s *Storage) GetPartially(ctx context.Context, path string, opt *storage.GetPartiallyRequestOptions) storage.IResponse {
	path = strings.TrimPrefix(path, "/")
	if opt == nil || len(opt.Range) == 0 {
		return s.Get(ctx, path)
	}
	offset, length, err := s.parseRange(opt.Range)
	if err != nil {
		return storage.NewResponse([]byte(nil), http.StatusUnprocessableEntity, ErrInvalidRange)
	}
	objHandle := s.bucketHandle.Object(path)
	r, err := objHandle.NewRangeReader(ctx, offset, length)
	if errors.Is(err, gs.ErrObjectNotExist) {
		return storage.NewResponse(nil, http.StatusNotFound, err)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		return storage.NewResponse(nil, apiErr.Code, apiErr)
	}
	defer r.Close()
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return storage.NewResponse(nil, http.StatusUnprocessableEntity, err)
	}
	objAttrs, err := objHandle.Attrs(ctx)
	if err != nil {
		return storage.NewResponse(nil, http.StatusNotFound, err)
	}
	return storage.NewResponse(d, http.StatusPartialContent, nil).
		WithMetadata(s.parseMetadata(objAttrs, offset, length))
}

func (s *Storage) parseRange(input string) (int64, int64, error) {
	matches := rangeRegex.FindStringSubmatch(input)
	if matches == nil {
		return 0, 0, errors.New("range parse error")
	}
	start, _ := strconv.ParseInt(matches[1], 10, 64)
	end, _ := strconv.ParseInt(matches[2], 10, 64)
	return start, (end - start) + 1, nil
}

func (s *Storage) parseMetadata(attrs *gs.ObjectAttrs, offset, length int64) *storage.ResponseMetadata {
	return &storage.ResponseMetadata{
		// TODO: bytes is the only range unit formally defined by RFC 7233,
		// update this when GCS supports getting it via headers.
		// Ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
		AcceptRanges:  "bytes",
		ContentLength: strconv.FormatInt(int64(math.Abs(float64(offset-length))), 10),
		ContentRange:  fmt.Sprintf("bytes %d-%d/%d", offset, length-1, attrs.Size),
		ContentType:   attrs.ContentType,
		ETag:          attrs.Etag,
		LastModified:  attrs.Updated.Format(time.RFC1123),
	}
}

func clientOptions(credentialsJSON []byte) []option.ClientOption {
	if len(credentialsJSON) != 0 {
		return []option.ClientOption{option.WithCredentialsJSON(credentialsJSON)}
	}
	return []option.ClientOption{option.WithoutAuthentication()}
}
