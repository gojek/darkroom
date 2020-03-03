package storage

import "github.com/afex/hystrix-go/hystrix"

// ResponseMetadata contains metadata of the storage response
type ResponseMetadata struct {
	AcceptRanges  string
	ContentLength string
	ContentRange  string
	ContentType   string
	ETag          string
	LastModified  string
}

// Response implements the IResponse interface
type Response struct {
	data     []byte
	err      error
	status   int
	metadata *ResponseMetadata
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

// Metadata returns metadata field from the struct
func (r *Response) Metadata() *ResponseMetadata {
	return r.metadata
}

// WithMetadata sets metadata field on the struct
func (r *Response) WithMetadata(metadata *ResponseMetadata) *Response {
	r.metadata = metadata
	return r
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

// GetPartiallyRequestOptions holds option to request data from storage
type GetPartiallyRequestOptions struct {
	Range string
}
