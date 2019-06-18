package storage

import "github.com/afex/hystrix-go/hystrix"

type Response struct {
	data   []byte
	err    error
	status int
}

func (r *Response) Data() []byte {
	return r.data
}

func (r *Response) Error() error {
	return r.err
}

func (r *Response) Status() int {
	return r.status
}

func NewResponse(data []byte, statusCode int, err error) *Response {
	return &Response{data: data, err: err, status: statusCode}
}

type HystrixCommand struct {
	Name   string
	Config hystrix.CommandConfig
}
