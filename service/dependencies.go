package service

import (
	"***REMOVED***/darkroom/server/storage"
	base "***REMOVED***/darkroom/storage"
)

type Dependencies struct {
	Storage base.Storage
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		Storage: storage.NewS3Storage(),
	}
}
