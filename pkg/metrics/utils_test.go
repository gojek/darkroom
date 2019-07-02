package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetImageSizeCluster(t *testing.T) {
	assert.Equal(t, "<=128KB", GetImageSizeCluster(make([]byte, 128*1024)), "<=500KB")
	assert.Equal(t, "<=256KB", GetImageSizeCluster(make([]byte, 256*1024)), "<=500KB")
	assert.Equal(t, "<=512KB", GetImageSizeCluster(make([]byte, 512*1024)), "<=500KB")
	assert.Equal(t, "<=1MB", GetImageSizeCluster(make([]byte, 1024*1024)), "<=500KB")
	assert.Equal(t, "<=2MB", GetImageSizeCluster(make([]byte, 2048*1024)), "<=500KB")
	assert.Equal(t, ">2MB", GetImageSizeCluster(make([]byte, 2049*1024)), "<=500KB")
}
