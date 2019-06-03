package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDependencies(t *testing.T) {
	deps := NewDependencies()
	assert.NotNil(t, deps)
}
