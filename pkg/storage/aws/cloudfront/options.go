package cloudfront

import (
	"github.com/gojektech/heimdall"
)

// Option represents the Cloudfront storage options
type Option func(s *Storage)

// WithCloudfrontHost sets the cloudfront host, can ends up with trailing slash or not
func WithCloudfrontHost(host string) Option {
	return func(s *Storage) {
		s.cloudfrontHost = host
	}
}

// WithHeimdallClient sets the client
func WithHeimdallClient(client heimdall.Client) Option {
	return func(s *Storage) {
		s.client = client
	}
}

// WithSecureProtocol uses https while making requests with the client
func WithSecureProtocol() Option {
	return func(s *Storage) {
		s.secureProtocol = true
	}
}
