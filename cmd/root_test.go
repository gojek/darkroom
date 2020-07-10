package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootCmd(t *testing.T) {
	// when
	err := Execute()
	// then
	assert.NoError(t, err)
}
