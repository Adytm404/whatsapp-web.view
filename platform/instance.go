package platform

// +build windows linux darwin freebsd

type InstanceLock interface {
	Acquire() bool
	Restore()
	Release()
}
