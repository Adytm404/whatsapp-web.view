package main

// +build windows linux darwin freebsd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"whatsapp-desktop/platform"
)

const (
	windowTitle = "WhatsApp Desktop"
	appURL      = "https://web.whatsapp.com"
	mutexName   = "WhatsAppDesktopSingleInstanceMutex"
	appDirName  = "WhatsAppDesktopLight"
)

func main() {
	lock := platform.NewInstanceLock(mutexName)
	if !lock.Acquire() {
		lock.Restore()
		os.Exit(0)
	}
	defer lock.Release()

	userDataDir := getUserDataDir()

	w := platform.NewWebView(platform.WebViewConfig{
		Title:    windowTitle,
		Width:    1100,
		Height:   750,
		DataPath: userDataDir,
		Debug:    false,
	})
	if w == nil {
		log.Fatalln("Gagal inisialisasi WebView")
	}
	defer w.Destroy()

	platform.SetDarkFrame(w.Window())

	notifier := platform.NewNotifier()
	_ = w.Bind("sendNativeNotification", func(title, body string) {
		go notifier.Push(title, body, "")
	})

	w.Init(getInitScript(platform.OSUserAgent))
	w.Navigate(appURL)
	w.Run()
}

func getUserDataDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	dir := filepath.Join(configDir, appDirName, "UserData")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func getInitScript(ua string) string {
	return fmt.Sprintf(`
		Object.defineProperty(navigator, 'userAgent', {
			get: () => 'Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36'
		});
		Object.defineProperty(navigator, 'appVersion', {
			get: () => 'Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36'
		});
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
		var darkCSS = document.createElement('style');
		darkCSS.textContent = 'headerbar, .titlebar { background-color: #111B21 !important; color: #FFFFFF !important; }';
		document.head.appendChild(darkCSS);
	`, ua, ua)
}
