// Package server holds the structs and methods containing the logic to spin up the api server
package server

// LifeCycleHook can be used to add custom functions in the server lifecycle
type LifeCycleHook struct {
	initFunc  func()
	deferFunc func()
}

// NewLifeCycleHook takes an initFunc and a deferFunc as arguments and returns a new LifeCycleHook
func NewLifeCycleHook(initFunc func(), deferFunc func()) *LifeCycleHook {
	return &LifeCycleHook{
		initFunc:  initFunc,
		deferFunc: deferFunc,
	}
}
