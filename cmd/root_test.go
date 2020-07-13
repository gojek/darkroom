package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootCmd(t *testing.T) {
	err := Execute()
	assert.NoError(t, err)
}
