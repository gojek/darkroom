package service

import (
	"***REMOVED***/darkroom/processor/native"
	"***REMOVED***/darkroom/server/storage"
	base "***REMOVED***/darkroom/storage"
)

type Dependencies struct {
	Storage     base.Storage
	Manipulator Manipulator
}

func NewDependencies() *Dependencies {
	return &Dependencies{
		Storage:     storage.NewS3Storage(),
		Manipulator: NewManipulator(native.NewBildProcessor()),
	}
}
