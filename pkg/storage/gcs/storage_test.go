package gcs

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	storageTypes "github.com/gojek/darkroom/pkg/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	validPath      = "path/to/valid-file"
	invalidPath    = "path/to/invalid-file"
	unreadablePath = "path/to/unreadable-file"
	validRange     = "bytes=0-100"
	invalidRange   = "none"
	bucketName     = "bucket-name"
)

type StorageTestSuite struct {
	suite.Suite
	storage      Storage
	bucketHandle *mockBucketHandle
}

func (s *StorageTestSuite) SetupTest() {
	ns, err := NewStorage(Options{BucketName: bucketName})
	s.NoError(err)
	s.bucketHandle = &mockBucketHandle{}
	ns.bucketHandle = s.bucketHandle
	s.storage = *ns
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (s *StorageTestSuite) TestNewStorage() {
	s.NotNil(s.storage)
}

func (s *StorageTestSuite) TestNewStorageWithValidCredentialsJSON() {
	ns, err := NewStorage(Options{BucketName: bucketName, CredentialsJSON: []byte(`
{
  "type": "service_account",
  "project_id": "",
  "private_key_id": "",
  "private_key": "",
  "client_email": "",
  "client_id": ""
}`)})
	s.NoError(err)
	s.NotNil(ns)
}

func (s *StorageTestSuite) TestNewStorageWithInvalidCredentialsJSON() {
	ns, err := NewStorage(Options{BucketName: bucketName, CredentialsJSON: []byte("randomJson")})
	s.Error(err)
	s.Nil(ns)
}

func (s *StorageTestSuite) TestNewStorageHasBucketHandle() {
	ns, err := NewStorage(Options{BucketName: bucketName})
	s.NoError(err)
	s.NotNil(ns)
	s.NotNil(ns.bucketHandle)
}

func (s *StorageTestSuite) TestNewStorageHasCorrectBucketName() {
	ns, err := NewStorage(Options{BucketName: bucketName})
	s.NoError(err)
	s.NotNil(ns)
	c, _ := storage.NewClient(
		context.TODO(),
		option.WithHTTPClient(newTestClient(bucketResponseMocker)),
	)
	ns.bucketHandle = bucketHandle{c.Bucket(bucketName)}
	attrs, err := ns.bucketHandle.Attrs(context.TODO())
	s.NoError(err)
	s.Equal(bucketName, attrs.Name)
}

func (s *StorageTestSuite) TestBenchForStorage_Get() {
	errForbidden := &googleapi.Error{Code: 403, Message: "Forbidden"}
	testcases := []struct {
		name            string
		ctx             context.Context
		path            string
		newReaderReturn func() (Reader, error)
		res             storageTypes.IResponse
	}{
		{
			name: "Success",
			ctx:  context.TODO(),
			path: validPath,
			newReaderReturn: func() (Reader, error) {
				return ioutil.NopCloser(strings.NewReader("someData")), nil
			},
			res: storageTypes.NewResponse([]byte("someData"), http.StatusOK, nil),
		},
		{
			name: "FailureWithInvalidPath",
			ctx:  context.TODO(),
			path: invalidPath,
			newReaderReturn: func() (Reader, error) {
				return nil, storage.ErrObjectNotExist
			},
			res: storageTypes.NewResponse([]byte(nil), http.StatusNotFound, storage.ErrObjectNotExist),
		},
		{
			name: "FailureWithForbiddenPath",
			ctx:  context.TODO(),
			path: validPath,
			newReaderReturn: func() (Reader, error) {
				return nil, errForbidden
			},
			res: storageTypes.NewResponse([]byte(nil), http.StatusForbidden, errForbidden),
		},
		{
			name: "FailureWithUnreadablePath",
			ctx:  context.TODO(),
			path: unreadablePath,
			newReaderReturn: func() (Reader, error) {
				return &badReader{}, nil
			},
			res: storageTypes.NewResponse([]byte(nil), http.StatusUnprocessableEntity, io.ErrUnexpectedEOF),
		},
	}

	for _, t := range testcases {
		s.SetupTest()
		s.Run(t.name, func() {
			mo := &mockObjectHandle{objectKey: t.path}
			s.bucketHandle.On("Object", t.path).Return(mo)
			mo.On("NewReader", t.ctx).Return(t.newReaderReturn())

			res := s.storage.Get(t.ctx, t.path)

			s.Equal(t.res.Error(), res.Error())
			s.Equal(t.res.Data(), res.Data())
			s.Equal(t.res.Status(), res.Status())
			s.Equal(t.res.Metadata(), res.Metadata())
		})
	}
}

func (s *StorageTestSuite) TestStorageRangeGetter() {
	offset, length, err := s.storage.parseRange("bytes=100-200")
	s.Equal(int64(100), offset)
	s.Equal(int64(101), length)
	s.NoError(err)

	offset, length, err = s.storage.parseRange(invalidRange)
	s.Equal(int64(0), offset)
	s.Equal(int64(0), length)
	s.Error(err)
}

func (s *StorageTestSuite) TestBenchForStorage_GetPartially() {
	validRange := validRange
	invalidRange := invalidRange
	outOfBoundRange := "bytes=4000-5000"
	emptyRange := ""
	testcases := []struct {
		name                 string
		ctx                  context.Context
		path                 string
		range_               *string
		newReaderReturn      func() (Reader, error)
		newRangeReaderReturn func() (Reader, error)
		attrsReturn          func() (*storage.ObjectAttrs, error)
		res                  storageTypes.IResponse
	}{
		{
			name:   "Success",
			ctx:    context.TODO(),
			path:   validPath,
			range_: &validRange,
			newRangeReaderReturn: func() (Reader, error) {
				return ioutil.NopCloser(strings.NewReader("someData")), nil
			},
			attrsReturn: func() (*storage.ObjectAttrs, error) {
				t, _ := time.Parse(time.RFC1123, "Wed, 21 Oct 2015 07:28:00 GMT")
				return &storage.ObjectAttrs{
					Bucket:      bucketName,
					Name:        validPath,
					ContentType: "image/png",
					Size:        247103,
					Updated:     t,
					Etag:        "32705ce195789d7bf07f3d44783c2988",
				}, nil
			},
			res: storageTypes.NewResponse([]byte("someData"), http.StatusPartialContent, nil).
				WithMetadata(&storageTypes.ResponseMetadata{
					AcceptRanges:  "bytes",
					ContentLength: "101",
					ContentRange:  "bytes 0-100/247103",
					ContentType:   "image/png",
					ETag:          "32705ce195789d7bf07f3d44783c2988",
					LastModified:  "Wed, 21 Oct 2015 07:28:00 GMT",
				}),
		},
		{
			name: "WithNilRequestOptions",
			ctx:  context.TODO(),
			path: validPath,
			newReaderReturn: func() (Reader, error) {
				return ioutil.NopCloser(strings.NewReader("someData")), nil
			},
			res: storageTypes.NewResponse([]byte("someData"), http.StatusOK, nil),
		},
		{
			name:   "WithEmptyRangeValue",
			ctx:    context.TODO(),
			path:   validPath,
			range_: &emptyRange,
			newReaderReturn: func() (Reader, error) {
				return ioutil.NopCloser(strings.NewReader("someData")), nil
			},
			res: storageTypes.NewResponse([]byte("someData"), http.StatusOK, nil),
		},
		{
			name: "OnInvalidPath",
			ctx:  context.TODO(),
			path: invalidPath,
			newReaderReturn: func() (Reader, error) {
				return nil, storage.ErrObjectNotExist
			},
			res: storageTypes.NewResponse(nil, http.StatusNotFound, storage.ErrObjectNotExist),
		},
		{
			name:   "OnInvalidRange",
			ctx:    context.TODO(),
			path:   validPath,
			range_: &invalidRange,
			newReaderReturn: func() (Reader, error) {
				return ioutil.NopCloser(strings.NewReader("someData")), nil
			},
			res: storageTypes.NewResponse(nil, http.StatusUnprocessableEntity, ErrInvalidRange),
		},
		{
			name:   "OnRangeReaderError",
			ctx:    context.TODO(),
			path:   validPath,
			range_: &outOfBoundRange,
			newRangeReaderReturn: func() (Reader, error) {
				return nil, &googleapi.Error{Code: 400, Message: "Bad Request"}
			},
			res: storageTypes.NewResponse(nil, http.StatusBadRequest, &googleapi.Error{
				Code:    400,
				Message: "Bad Request"},
			),
		},
		{
			name:   "NotFoundWithValidRangeAndInvalidPath",
			ctx:    context.TODO(),
			path:   invalidPath,
			range_: &validRange,
			newRangeReaderReturn: func() (Reader, error) {
				return nil, storage.ErrObjectNotExist
			},
			res: storageTypes.NewResponse([]byte(nil), http.StatusNotFound, storage.ErrObjectNotExist),
		},
		{
			name:   "OnBadReaderError",
			ctx:    context.TODO(),
			path:   unreadablePath,
			range_: &validRange,
			newRangeReaderReturn: func() (Reader, error) {
				return &badReader{}, nil
			},
			res: storageTypes.NewResponse(nil, http.StatusUnprocessableEntity, io.ErrUnexpectedEOF),
		},
		{
			name:   "WhenObjectAttributesAreUnreadable",
			ctx:    context.TODO(),
			path:   invalidPath,
			range_: &validRange,
			newRangeReaderReturn: func() (Reader, error) {
				return ioutil.NopCloser(strings.NewReader("someData")), nil
			},
			attrsReturn: func() (*storage.ObjectAttrs, error) {
				return nil, storage.ErrObjectNotExist
			},
			res: storageTypes.NewResponse(nil, http.StatusNotFound, storage.ErrObjectNotExist),
		},
	}

	for _, t := range testcases {
		s.SetupTest()
		s.Run(t.name, func() {
			mo := &mockObjectHandle{objectKey: t.path}
			s.bucketHandle.On("Object", t.path).Return(mo)

			var opts *storageTypes.GetPartiallyRequestOptions
			if t.range_ != nil {
				if o, l, err := s.storage.parseRange(*t.range_); err == nil {
					mo.On("NewRangeReader", t.ctx, o, l).
						Return(t.newRangeReaderReturn())
				}
				opts = &storageTypes.GetPartiallyRequestOptions{Range: *t.range_}
			}
			if t.newReaderReturn != nil {
				mo.On("NewReader", t.ctx).Return(t.newReaderReturn())
			}
			if t.attrsReturn != nil {
				mo.On("Attrs", t.ctx).Return(t.attrsReturn())
			}

			res := s.storage.GetPartially(t.ctx, t.path, opts)

			s.Equal(t.res.Error(), res.Error())
			s.Equal(t.res.Data(), res.Data())
			s.Equal(t.res.Status(), res.Status())
			s.Equal(t.res.Metadata(), res.Metadata())
		})
	}
}

type mockBucketHandle struct {
	mock.Mock
}

func (m *mockBucketHandle) Object(s string) ObjectHandle {
	args := m.Called(s)
	return args[0].(*mockObjectHandle)
}

func (m *mockBucketHandle) Attrs(ctx context.Context) (*storage.BucketAttrs, error) {
	args := m.Called(ctx)
	return args[0].(*storage.BucketAttrs), args.Error(1)
}

type mockObjectHandle struct {
	mock.Mock
	objectKey string
}

func (m *mockObjectHandle) Attrs(ctx context.Context) (attrs *storage.ObjectAttrs, err error) {
	args := m.Called(ctx)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*storage.ObjectAttrs), args.Error(1)
}

func (m *mockObjectHandle) NewReader(ctx context.Context) (Reader, error) {
	args := m.Called(ctx)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(Reader), args.Error(1)
}

func (m *mockObjectHandle) NewRangeReader(ctx context.Context, offset, length int64) (Reader, error) {
	args := m.Called(ctx, offset, length)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(Reader), args.Error(1)
}

type badReader struct{}

func (b badReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (b badReader) Close() error {
	return nil
}
