package service

import (
	"***REMOVED***/darkroom/processor/native"
	"***REMOVED***/darkroom/server/config"
	base "***REMOVED***/darkroom/storage"
	"***REMOVED***/darkroom/storage/s3"
)

type Dependencies struct {
	Storage     base.Storage
	Manipulator Manipulator
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		Storage: s3.NewStorage(
			s3.WithBucketName(config.BucketName()),
			s3.WithBucketRegion(config.BucketRegion()),
			s3.WithAccessKey(config.BucketAccessKey()),
			s3.WithSecretKey(config.BucketSecretKey()),
			s3.WithHystrixCommand(config.HystrixCommand()),
		),
		Manipulator: NewManipulator(native.NewBildProcessor()),
	}
}
