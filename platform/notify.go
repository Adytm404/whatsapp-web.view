//go:build windows || linux || darwin || freebsd

package platform

type Notifier interface {
	Push(title, message, iconPath string) error
}
