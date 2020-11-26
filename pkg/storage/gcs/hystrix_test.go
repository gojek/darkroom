package gcs

import (
	"context"
	"net/http"
	"testing"

	"github.com/gojektech/heimdall/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestNewHeimdallHTTPClientWithInvalidCredentials(t *testing.T) {
	hc := hystrix.NewClient()
	hhc, err := newHeimdallHTTPClient(context.TODO(), hc, []byte("random"))
	assert.Nil(t, hhc)
	assert.Error(t, err)
}

func TestNewHeimdallHTTPClientWithNoCredentials(t *testing.T) {
	hc := hystrix.NewClient()
	hhc, err := newHeimdallHTTPClient(context.TODO(), hc, []byte(""))
	assert.NotNil(t, hhc)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	_, err = hhc.Do(req)
	assert.Error(t, err, "expecting unsupported protocol error")
}
