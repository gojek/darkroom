package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecBuilder_Build(t *testing.T) {
	scope := "scope"
	spec := NewSpecBuilder().WithScope(scope).Build()
	assert.Equal(t, spec.Scope, scope)
}
