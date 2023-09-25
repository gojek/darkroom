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
	hhc, err := newHeimdallHTTPClient(context.TODO(), &Options{
		CredentialsJSON: []byte("random"),
		Client:          hc,
	})
	assert.Nil(t, hhc)
	assert.Error(t, err)
}

func TestNewHeimdallHTTPClientWithNoCredentials(t *testing.T) {
	hc := hystrix.NewClient()
	hhc, err := newHeimdallHTTPClient(context.TODO(), &Options{
		CredentialsJSON: []byte(""),
		Client:          hc,
	})
	assert.NotNil(t, hhc)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	_, err = hhc.Do(req)
	assert.Error(t, err, "expecting unsupported protocol error")
}
