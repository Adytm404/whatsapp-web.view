//go:build windows

package platform

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	dwmapi         = windows.NewLazySystemDLL("dwmapi.dll")
	procDwmSetAttr = dwmapi.NewProc("DwmSetWindowAttribute")
)

const (
	DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1 = 19
	DWMWA_USE_IMMERSIVE_DARK_MODE             = 20
	DWMWA_CAPTION_COLOR                      = 35
	DWMWA_TEXT_COLOR                         = 36
)

func SetDarkFrame(window uintptr) {
	darkMode := int32(1)
	procDwmSetAttr.Call(
		window,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)
	procDwmSetAttr.Call(
		window,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)
	captionColor := uint32(0x00211B11)
	procDwmSetAttr.Call(
		window,
		uintptr(DWMWA_CAPTION_COLOR),
		uintptr(unsafe.Pointer(&captionColor)),
		unsafe.Sizeof(captionColor),
	)
	textColor := uint32(0x00FFFFFF)
	procDwmSetAttr.Call(
		window,
		uintptr(DWMWA_TEXT_COLOR),
		uintptr(unsafe.Pointer(&textColor)),
		unsafe.Sizeof(textColor),
	)
}
