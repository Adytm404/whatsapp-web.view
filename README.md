# WhatsApp Webview Desktop

Lightweight cross-platform desktop wrapper for [WhatsApp Web](https://web.whatsapp.com), built with Go and system WebView.

## Features

- Native WebView window (no Electron)
- Persistent WhatsApp session across app restarts
- Cookies, LocalStorage, IndexedDB, and service-worker data stored in a dedicated profile
- Dark title bar theme
- Native notification bridge
- Camera and microphone access for WhatsApp voice and video calls
- Single-instance protection
- High-DPI display support
- Small native executable

## Supported Platforms

| Platform | WebView Engine | Status |
|----------|---------------|--------|
| Windows 10/11 | Microsoft Edge WebView2 | ✅ Supported |
| Ubuntu/Debian | WebKitGTK | ✅ Supported |
| Arch Linux | WebKitGTK | ✅ Supported |
| macOS | WKWebView | 🔜 Planned |
| FreeBSD | WebKitGTK | 🔜 Planned |

## Requirements

### Windows
- Windows 10 or newer
- Microsoft Edge WebView2 Runtime (auto-downloaded if missing)
- WhatsApp account paired with WhatsApp Web

### Linux (Ubuntu/Debian)
```bash
sudo apt install libwebkit2gtk-4.1-0 libgtk-3-0
```

### Linux (Arch Linux)
```bash
sudo pacman -S webkit2gtk-4.1 gtk3
```

## Download

### Windows
Download `WhatsApp.exe` from the latest [GitHub Release](https://github.com/sandikodev/whatsapp-webview/releases).

### Ubuntu/Debian
```bash
sudo dpkg -i whatsapp-webview_1.0.0_amd64.deb
sudo apt install -f
```

### Arch Linux (AUR)
```bash
paru -S whatsapp-webview-git
```

## Build From Source

### Prerequisites
- Go 1.22 or newer

### Build for Linux
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o whatsapp-webview .
```

### Build for Windows
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .
```

### Build for macOS
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o whatsapp-webview .
```

## Session Data

Profile data is stored at:
- **Windows**: `%APPDATA%\WhatsAppDesktopLight\UserData`
- **Linux**: `~/.config/WhatsAppDesktopLight/UserData`
- **macOS**: `~/Library/Application Support/WhatsAppDesktopLight/UserData`

Do not delete this folder if the existing login session must remain available.

## Project Structure

```
├── main.go                    # Shared entry point
├── platform/
│   ├── webview.go             # WebView interface
│   ├── notify.go              # Notifier interface
│   ├── instance.go            # InstanceLock interface
│   ├── webview_linux.go       # Linux WebView (glaze/WebKitGTK)
│   ├── webview_windows.go     # Windows WebView (go-webview2)
│   ├── notify_linux.go        # Linux notifications (D-Bus)
│   ├── notify_windows.go      # Windows notifications (Toast)
│   ├── instance_linux.go      # Linux instance lock (flock)
│   ├── instance_windows.go    # Windows instance lock (Mutex)
│   ├── frame_linux.go         # Linux dark mode
│   ├── frame_windows.go       # Windows dark mode (DWM)
│   └── ua_*.go                # User-Agent per platform
├── packaging/
│   ├── ubuntu/                # Ubuntu .deb package
│   └── archlinux/             # Arch Linux PKGBUILD
└── assets/
    └── icon.*                 # App icons
```

## Privacy

This app loads WhatsApp Web directly. Chat data and authentication state are handled by WhatsApp Web and stored locally in the WebView profile above. This project is not affiliated with WhatsApp or Meta.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Credits

- Original Windows-only project: [Adytm404/whatsapp-web.view](https://github.com/Adytm404/whatsapp-web.view)
- Cross-platform WebView: [crgimenes/glaze](https://github.com/crgimenes/glaze)
