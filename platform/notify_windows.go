//go:build windows

package platform

import "github.com/go-toast/toast"

type windowsNotifier struct{}

func NewNotifier() Notifier {
	return &windowsNotifier{}
}

func (n *windowsNotifier) Push(title, message, iconPath string) error {
	notification := toast.Notification{
		AppID:   "WhatsApp Desktop",
		Title:   title,
		Message: message,
		Icon:    iconPath,
	}
	return notification.Push()
}
