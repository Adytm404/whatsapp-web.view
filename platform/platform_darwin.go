//go:build darwin

package platform

import "fmt"

type darwinInstanceLock struct{}

func NewInstanceLock(name string) InstanceLock {
	return &darwinInstanceLock{}
}

func (l *darwinInstanceLock) Acquire() bool           { return true }
func (l *darwinInstanceLock) Restore()                {}
func (l *darwinInstanceLock) Release()                {}

type darwinWebView struct{}

func NewWebView(cfg WebViewConfig) WebView {
	fmt.Println("macOS WebView: glaze not yet implemented")
	return &darwinWebView{}
}

func (v *darwinWebView) SetTitle(title string)                {}
func (v *darwinWebView) SetSize(width, height int, hint ResizeHint) {}
func (v *darwinWebView) Navigate(url string)                  {}
func (v *darwinWebView) Init(script string)                   {}
func (v *darwinWebView) Bind(name string, fn interface{}) error { return nil }
func (v *darwinWebView) Run()                                 {}
func (v *darwinWebView) Destroy()                             {}
func (v *darwinWebView) Window() uintptr                      { return 0 }

type darwinNotifier struct{}

func NewNotifier() Notifier {
	return &darwinNotifier{}
}

func (n *darwinNotifier) Push(title, message, iconPath string) error {
	fmt.Printf("macOS notification: %s - %s\n", title, message)
	return nil
}

func SetDarkFrame(window uintptr) {}

const OSUserAgent = "Macintosh; Intel Mac OS X 14_0"
