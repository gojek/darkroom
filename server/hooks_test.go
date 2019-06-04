package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLifeCycleHook(t *testing.T) {
	initCount := 0
	deferCount := 0

	func1 := func() { initCount++ }
	func2 := func() { deferCount++ }

	hook := NewLifeCycleHook(func1, func2)

	hook.initFunc()
	hook.deferFunc()

	assert.Equal(t, 1, initCount)
	assert.Equal(t, 1, deferCount)
}
