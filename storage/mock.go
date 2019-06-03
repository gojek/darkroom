package storage

import (
	"context"
	"github.com/stretchr/testify/mock"
	"***REMOVED***/darkroom/storage"
)

type MockBaseStorage struct {
	mock.Mock
}

func (m *MockBaseStorage) Get(ctx context.Context, path string) storage.IResponse {
	args := m.Called(ctx, path)
	return storage.NewResponse(args[0].([]byte), args.Int(1), args.Error(2))
}

