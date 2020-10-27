package gcs

import (
	"net/http"
	"testing"

	"github.com/gojektech/heimdall/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestNewHeimdallHTTPClient(t *testing.T) {
	hc := hystrix.NewClient()
	hhc := newHeimdallHTTPClient(hc)
	assert.NotNil(t, hhc)
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	_, err := hhc.Do(req)
	assert.Error(t, err, "expecting unsupported protocol error")
}
