package cloudfront

import (
	"testing"

	"github.com/gojektech/heimdall/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestWithCloudfrontHost(t *testing.T) {
	hc := hystrix.NewClient()
	s := NewStorage(
		WithCloudfrontHost("cloudfront.net"),
		WithHeimdallClient(hc),
		WithSecureProtocol(),
	)
	assert.Equal(t, "cloudfront.net", s.cloudfrontHost)
	assert.Equal(t, hc, s.client)
	assert.Equal(t, true, s.secureProtocol)
}
