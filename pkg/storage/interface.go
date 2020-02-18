// Package storage contains the Storage interface and its various implementations for different backends.
package storage

import "context"

// Storage interface sets the contract that the implementation has to fulfil.
type Storage interface {
	// Get takes in the Context and path as an argument and returns an IResponse interface implementation.
	// This method figures out how to get the data from the storage backend.
	Get(ctx context.Context, path string, opt *GetRequestOptions) IResponse
}

// IResponse interface sets the contract that can be used to return the result from different Storage backends.
type IResponse interface {
	// Data method returns a byte array if the operation was successful
	Data() []byte
	// Error method returns an error if the operation was unsuccessful
	Error() error
	// Status method returns the http StatusCode from the storage backend if there is any
	Status() int
	// Metadata method returns response metadata if the operation is successful and client support metadata
	Metadata() *ResponseMetadata
}
