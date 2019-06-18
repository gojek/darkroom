package storage

import "context"

type Storage interface {
	Get(ctx context.Context, path string) IResponse
}

type IResponse interface {
	Data() []byte
	Error() error
	Status() int
}
