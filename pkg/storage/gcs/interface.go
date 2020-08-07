package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
)

// These interfaces allow us to test the package properly.
// See https://github.com/googleapis/google-cloud-go-testing/blob/master/storage/stiface/interfaces.go

type ObjectHandle interface {
	NewReader(context.Context) (Reader, error)
	NewRangeReader(context.Context, int64, int64) (Reader, error)
}

type BucketHandle interface {
	Object(string) ObjectHandle
}

type Reader interface {
	io.ReadCloser
}

type (
	bucketHandle struct{ *storage.BucketHandle }
	objectHandle struct{ *storage.ObjectHandle }
	reader       struct{ *storage.Reader }
)

func (b bucketHandle) Object(name string) ObjectHandle {
	return objectHandle{b.BucketHandle.Object(name)}
}

func (o objectHandle) NewReader(ctx context.Context) (Reader, error) {
	r, err := o.ObjectHandle.NewReader(ctx)
	return reader{r}, err
}

func (o objectHandle) NewRangeReader(ctx context.Context, offset, length int64) (Reader, error) {
	r, err := o.ObjectHandle.NewRangeReader(ctx, offset, length)
	return reader{r}, err
}
