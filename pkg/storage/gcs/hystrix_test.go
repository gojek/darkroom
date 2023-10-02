package gcs

import (
	"context"
	"golang.org/x/oauth2/google"
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

func TestNewHeimdallHTTPClientWithCustomCredentials(t *testing.T) {
	hc := hystrix.NewClient()
	hhc, err := newHeimdallHTTPClient(context.TODO(), &Options{
		CredentialsJSON: nil,
		Credentials:     &google.Credentials{ProjectID: "sample"},
		Client:          hc,
	})
	assert.NotNil(t, hhc)
	assert.NoError(t, err)
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	_, err = hhc.Do(req)
	assert.Error(t, err, "expecting unsupported protocol error")
}
