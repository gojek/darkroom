package storage

import "github.com/afex/hystrix-go/hystrix"

// Response implements the IResponse interface
type Response struct {
	data   []byte
	err    error
	status int
}

// Data returns the data field from the struct
func (r *Response) Data() []byte {
	return r.data
}

// Error returns the err field from the struct
func (r *Response) Error() error {
	return r.err
}

// Status returns the status field from the struct
func (r *Response) Status() int {
	return r.status
}

// NewResponse takes data, statusCode and error as arguments and returns a new Response
func NewResponse(data []byte, statusCode int, err error) *Response {
	return &Response{data: data, err: err, status: statusCode}
}

// HystrixCommand wraps the command name and the configuration to be used with hystrix
type HystrixCommand struct {
	Name   string
	Config hystrix.CommandConfig
}
