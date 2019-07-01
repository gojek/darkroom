package s3

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gojek/darkroom/pkg/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionsAreSet(t *testing.T) {
	hystrixCmd := storage.HystrixCommand{
		Name: "TestCommand",
		Config: hystrix.CommandConfig{
			Timeout:                5000,
			MaxConcurrentRequests:  100,
			RequestVolumeThreshold: 10,
			SleepWindow:            10,
			ErrorPercentThreshold:  25,
		},
	}
	s := NewStorage(
		WithBucketName("bucket"),
		WithBucketRegion("region"),
		WithAccessKey("accessKey"),
		WithSecretKey("secretKey"),
		WithHystrixCommand(hystrixCmd),
	)

	assert.Equal(t, "bucket", s.bucketName)
	assert.Equal(t, "region", s.bucketRegion)
	assert.Equal(t, "accessKey", s.accessKey)
	assert.Equal(t, "secretKey", s.secretKey)
	assert.Equal(t, hystrixCmd, s.hystrixCmd)
}
