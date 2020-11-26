package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	err := Execute()
	assert.NoError(t, err)
}
