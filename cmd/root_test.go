package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootCmd(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{})
	// when
	err := cmd.Execute()
	// then
	assert.NoError(t, err)
}
