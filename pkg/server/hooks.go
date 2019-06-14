package server

type LifeCycleHook struct {
	initFunc  func()
	deferFunc func()
}

func NewLifeCycleHook(initFunc func(), deferFunc func()) *LifeCycleHook {
	return &LifeCycleHook{
		initFunc:  initFunc,
		deferFunc: deferFunc,
	}
}
