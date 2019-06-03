package storage

import (
	"context"
	"***REMOVED***/darkroom/server/config"
	"***REMOVED***/darkroom/storage"
	"***REMOVED***/darkroom/storage/s3"
)

type S3Storage struct {
	base       storage.Storage
	pathPrefix string
}

func (s3s S3Storage) Get(ctx context.Context, path string) storage.IResponse {
	if s3s.pathPrefix != "" {
		return s3s.base.Get(ctx, s3s.pathPrefix+path)
	} else {
		return s3s.base.Get(ctx, path)
	}
}

func NewS3Storage() storage.Storage {
	return S3Storage{
		base:       s3.NewStorage(config.BucketName(), config.BucketAccessKey(), config.BucketSecretKey()),
		pathPrefix: config.BucketPathPrefix(),
	}
}
