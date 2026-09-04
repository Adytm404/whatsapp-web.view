//go:build windows || linux || darwin || freebsd

package platform

type InstanceLock interface {
	Acquire() bool
	Restore()
	Release()
}
