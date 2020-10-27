package webfolder

import (
	"testing"

	"github.com/gojektech/heimdall/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestWithCloudfrontHost(t *testing.T) {
	hc := hystrix.NewClient()
	s := NewStorage(
		WithBaseURL("https://example.com/path/to/images"),
		WithHeimdallClient(hc),
	)
	assert.Equal(t, "https://example.com/path/to/images", s.baseURL)
	assert.Equal(t, hc, s.client)
}
