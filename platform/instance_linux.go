//go:build linux

package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type linuxInstanceLock struct {
	file *os.File
	path string
}

func NewInstanceLock(name string) InstanceLock {
	lockDir, _ := os.UserConfigDir()
	lockPath := filepath.Join(lockDir, "whatsapp-webview", "lock")
	os.MkdirAll(filepath.Dir(lockPath), 0755)
	return &linuxInstanceLock{path: lockPath}
}

func (l *linuxInstanceLock) Acquire() bool {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return true
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		f.Close()
		return false
	}
	l.file = f
	f.Truncate(0)
	f.Seek(0, 0)
	f.WriteString(fmt.Sprintf("%d", os.Getpid()))
	return true
}

func (l *linuxInstanceLock) Restore() {}

func (l *linuxInstanceLock) Release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.path)
	}
}
