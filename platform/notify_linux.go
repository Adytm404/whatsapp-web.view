package platform

// +build linux

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

type linuxNotifier struct {
	conn *dbus.Conn
}

func NewNotifier() Notifier {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Printf("D-Bus connection failed: %v\n", err)
		return &linuxNotifier{conn: nil}
	}
	return &linuxNotifier{conn: conn}
}

func (n *linuxNotifier) Push(title, message, iconPath string) error {
	if n.conn == nil {
		fmt.Printf("Notification: %s - %s\n", title, message)
		return nil
	}
	obj := n.conn.Object(
		"org.freedesktop.Notifications",
		"/org/freedesktop/Notifications",
	)
	call := obj.Call(
		"org.freedesktop.Notifications.Notify",
		0,
		"WhatsApp Desktop",
		uint32(0),
		"",
		title,
		message,
		[]string{},
		map[string]dbus.Variant{},
		int32(-1),
	)
	return call.Err
}
