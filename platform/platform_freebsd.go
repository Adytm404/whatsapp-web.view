//go:build freebsd

package platform

import "fmt"

type freebsdInstanceLock struct{}

func NewInstanceLock(name string) InstanceLock {
	return &freebsdInstanceLock{}
}

func (l *freebsdInstanceLock) Acquire() bool           { return true }
func (l *freebsdInstanceLock) Restore()                {}
func (l *freebsdInstanceLock) Release()                {}

type freebsdWebView struct{}

func NewWebView(cfg WebViewConfig) WebView {
	fmt.Println("FreeBSD WebView: glaze not yet implemented")
	return &freebsdWebView{}
}

func (v *freebsdWebView) SetTitle(title string)                {}
func (v *freebsdWebView) SetSize(width, height int, hint ResizeHint) {}
func (v *freebsdWebView) Navigate(url string)                  {}
func (v *freebsdWebView) Init(script string)                   {}
func (v *freebsdWebView) Bind(name string, fn interface{}) error { return nil }
func (v *freebsdWebView) Run()                                 {}
func (v *freebsdWebView) Destroy()                             {}
func (v *freebsdWebView) Window() uintptr                      { return 0 }

type freebsdNotifier struct{}

func NewNotifier() Notifier {
	return &freebsdNotifier{}
}

func (n *freebsdNotifier) Push(title, message, iconPath string) error {
	fmt.Printf("FreeBSD notification: %s - %s\n", title, message)
	return nil
}

func SetDarkFrame(window uintptr) {}

const OSUserAgent = "X11; FreeBSD amd64"
