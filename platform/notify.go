package platform

// +build windows linux darwin freebsd

type Notifier interface {
	Push(title, message, iconPath string) error
}
