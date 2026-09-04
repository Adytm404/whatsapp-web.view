//go:build windows

package platform

import (
	"syscall"
	"unsafe"

	"github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

var (
	kernel32        = windows.NewLazySystemDLL("kernel32.dll")
	user32          = windows.NewLazySystemDLL("user32.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
	procFindWindow  = user32.NewProc("FindWindowW")
	procSetFgWindow = user32.NewProc("SetForegroundWindow")
	procShowNormal  = user32.NewProc("ShowWindow")
)

type windowsInstanceLock struct {
	handle uintptr
}

func NewInstanceLock(name string) InstanceLock {
	return &windowsInstanceLock{}
}

func (l *windowsInstanceLock) Acquire() bool {
	namePtr, _ := syscall.UTF16PtrFromString("WhatsAppDesktopSingleInstanceMutex")
	handle, _, err := procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(namePtr)))
	l.handle = handle
	return err != windows.ERROR_ALREADY_EXISTS
}

func (l *windowsInstanceLock) Restore() {
	titlePtr, _ := syscall.UTF16PtrFromString("WhatsApp Desktop")
	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd != 0 {
		procShowNormal.Call(hwnd, 9) // SW_RESTORE
		procSetFgWindow.Call(hwnd)
	}
}

func (l *windowsInstanceLock) Release() {
	if l.handle != 0 {
		windows.CloseHandle(windows.Handle(l.handle))
	}
}

type windowsWebView struct {
	w webview2.WebView
}

func NewWebView(cfg WebViewConfig) WebView {
	opts := webview2.WebViewOptions{
		Window:    nil,
		Debug:     cfg.Debug,
		DataPath:  cfg.DataPath,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  cfg.Title,
			Width:  uint(cfg.Width),
			Height: uint(cfg.Height),
			IconId: 2,
			Center: true,
		},
	}
	w := webview2.NewWithOptions(opts)
	if w == nil {
		return nil
	}
	return &windowsWebView{w: w}
}

func (v *windowsWebView) SetTitle(title string) {
	v.w.SetTitle(title)
}

func (v *windowsWebView) SetSize(width, height int, hint ResizeHint) {
	v.w.SetSize(width, height, webview2.HintNone)
}

func (v *windowsWebView) Navigate(url string) {
	v.w.Navigate(url)
}

func (v *windowsWebView) Init(script string) {
	v.w.Init(script)
}

func (v *windowsWebView) Bind(name string, fn interface{}) error {
	return v.w.Bind(name, fn)
}

func (v *windowsWebView) Run() {
	v.w.Run()
}

func (v *windowsWebView) Destroy() {
	v.w.Destroy()
}

func (v *windowsWebView) Window() uintptr {
	return uintptr(v.w.Window())
}
