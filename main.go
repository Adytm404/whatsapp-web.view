package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/go-toast/toast"
	"github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

var (
	kernel32        = windows.NewLazySystemDLL("kernel32.dll")
	user32          = windows.NewLazySystemDLL("user32.dll")
	dwmapi          = windows.NewLazySystemDLL("dwmapi.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
	procFindWindow  = user32.NewProc("FindWindowW")
	procSetFgWindow = user32.NewProc("SetForegroundWindow")
	procShowNormal  = user32.NewProc("ShowWindow")
	procDwmSetAttr  = dwmapi.NewProc("DwmSetWindowAttribute")
)

const (
	windowTitle = "WhatsApp Desktop"
	appURL      = "https://web.whatsapp.com"
	mutexName   = "WhatsAppDesktopSingleInstanceMutex"
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"

	// DWM Window Attributes for Dark Theme
	DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1 = 19
	DWMWA_USE_IMMERSIVE_DARK_MODE             = 20
	DWMWA_CAPTION_COLOR                       = 35
	DWMWA_TEXT_COLOR                          = 36
)

func setDarkWindowFrame(hwnd uintptr) {
	darkMode := int32(1)
	// Try standard DWMWA_USE_IMMERSIVE_DARK_MODE (Win10 20H1+ & Win11)
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)
	// Try older Win10 build
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)

	// Set dark caption color (COLORREF: 0x00111B21 WhatsApp Dark Header: RGB 17, 27, 33)
	captionColor := uint32(0x00211B11) // 0x00BBGGRR
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_CAPTION_COLOR),
		uintptr(unsafe.Pointer(&captionColor)),
		unsafe.Sizeof(captionColor),
	)

	// Set white caption text (RGB 255, 255, 255)
	textColor := uint32(0x00FFFFFF)
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_TEXT_COLOR),
		uintptr(unsafe.Pointer(&textColor)),
		unsafe.Sizeof(textColor),
	)
}

func checkSingleInstance() (uintptr, bool) {
	namePtr, _ := syscall.UTF16PtrFromString(mutexName)
	handle, _, err := procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(namePtr)))
	if err == windows.ERROR_ALREADY_EXISTS {
		titlePtr, _ := syscall.UTF16PtrFromString(windowTitle)
		hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
		if hwnd != 0 {
			procShowNormal.Call(hwnd, 9) // SW_RESTORE
			procSetFgWindow.Call(hwnd)
		}
		return handle, false
	}
	return handle, true
}

func getUserDataDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = "."
		}
	}
	dir := filepath.Join(configDir, "WhatsAppDesktopLight", "UserData")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func showNativeNotification(title, message, iconPath string) {
	notification := toast.Notification{
		AppID:   "WhatsApp Desktop",
		Title:   title,
		Message: message,
		Icon:    iconPath,
	}
	_ = notification.Push()
}

func main() {
	_, isSingle := checkSingleInstance()
	if !isSingle {
		os.Exit(0)
	}
	userDataDir := getUserDataDir()
	executablePath, _ := os.Executable()
	iconFullPath := filepath.Join(filepath.Dir(executablePath), "icon.ico")

	opts := webview2.WebViewOptions{
		Window:    nil,
		Debug:     false,
		DataPath:  userDataDir,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  windowTitle,
			Width:  1100,
			Height: 750,
			IconId: 2,
			Center: true,
		},
	}

	w := webview2.NewWithOptions(opts)
	if w == nil {
		log.Fatalln("Gagal inisialisasi WebView2")
	}
	defer w.Destroy()

	hwnd := uintptr(w.Window())
	setDarkWindowFrame(hwnd)

	w.SetTitle(windowTitle)
	w.SetSize(1100, 750, webview2.HintNone)

	// Bind native notification bridge
	_ = w.Bind("sendNativeNotification", func(title, body string) {
		go showNativeNotification(title, body, iconFullPath)
	})

	// Inject JS: User-Agent spoofing + Notification API polyfill connecting to Go native Toast
	initScript := `
		// UserAgent override
		Object.defineProperty(navigator, 'userAgent', {
			get: () => '` + userAgent + `'
		});
		Object.defineProperty(navigator, 'appVersion', {
			get: () => '` + userAgent + `'
		});

		// Native Notification Polyfill for Windows Desktop Toast
		(function() {
			window.Notification = function(title, options) {
				options = options || {};
				var body = options.body || '';
				if (window.sendNativeNotification) {
					window.sendNativeNotification(title, body);
				}
				this.title = title;
				this.onclick = null;
				this.onclose = null;
				this.onerror = null;
				this.onshow = null;
			};
			window.Notification.permission = 'granted';
			window.Notification.requestPermission = function(callback) {
				var p = Promise.resolve('granted');
				if (typeof callback === 'function') {
					callback('granted');
				}
				return p;
			};
		})();

		// Persist zoom level in the WhatsApp Web profile.
		(function() {
			var storageKey = 'whatsapp-desktop-zoom';
			var zoom = parseFloat(localStorage.getItem(storageKey) || '1');
			if (!isFinite(zoom) || zoom < 0.5 || zoom > 2) zoom = 1;
			function saveZoom() {
				localStorage.setItem(storageKey, String(zoom));
				document.documentElement.style.zoom = String(zoom);
			}
			document.addEventListener('keydown', function(event) {
				if (!event.ctrlKey || event.altKey || event.metaKey) return;
				if (event.key === '+' || event.key === '=') {
					zoom = Math.min(2, Math.round((zoom + 0.1) * 10) / 10);
					saveZoom();
					event.preventDefault();
				} else if (event.key === '-' || event.key === '_') {
					zoom = Math.max(0.5, Math.round((zoom - 0.1) * 10) / 10);
					saveZoom();
					event.preventDefault();
				} else if (event.key === '0') {
					zoom = 1;
					saveZoom();
					event.preventDefault();
				}
			}, true);
			document.documentElement.style.zoom = String(zoom);
		})();
	`

	w.Init(initScript)
	w.Navigate(appURL)
	w.Run()
}
