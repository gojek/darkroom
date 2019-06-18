package webfolder

import (
	"github.com/gojektech/heimdall"
)

// Option represents the Cloudfront storage options
type Option func(s *Storage)

// WithBaseURL sets the baseURL
func WithBaseURL(url string) Option {
	return func(s *Storage) {
		s.baseURL = url
	}
}

// WithHeimdallClient sets the client
func WithHeimdallClient(client heimdall.Client) Option {
	return func(s *Storage) {
		s.client = client
	}
}
