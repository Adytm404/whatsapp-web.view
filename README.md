# WhatsApp Webview Desktop

Lightweight Windows desktop wrapper for [WhatsApp Web](https://web.whatsapp.com), built with Go and Microsoft Edge WebView2.

## Features

- Native Windows WebView2 window instead of Electron.
- Persistent WhatsApp session across app restarts.
- Cookies, LocalStorage, IndexedDB, and service-worker data stored in a dedicated profile.
- Dark Windows title bar and frame.
- Windows toast notification bridge.
- Camera and microphone access for WhatsApp voice and video calls.
- Persistent zoom level with `Ctrl +`, `Ctrl -`, and `Ctrl 0`.
- Standard clipboard shortcuts: `Ctrl C`, `Ctrl X`, `Ctrl V`, and `Ctrl A`.
- Single-instance protection.
- High-DPI display support.
- Small native executable.

## Requirements

- Windows 10 or newer.
- Microsoft Edge WebView2 Runtime. The app can download it automatically when missing.
- WhatsApp account paired with WhatsApp Web.

## Download

Download `WhatsApp.exe` from the latest [GitHub Release](https://github.com/Adytm404/whatsapp-web.view/releases).

Run the executable, scan the QR code, allow camera and microphone access when prompted by Windows, and keep using WhatsApp normally. The session is saved automatically.

## Session Data

Profile data is stored at:

```text
%APPDATA%\WhatsAppDesktopLight\UserData
```

Do not delete this folder if the existing login session must remain available. Closing the app does not clear session data.

## Build From Source

Install Go and a Windows C compiler, then run:

```powershell
go mod download
go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .
```

The repository includes the Windows manifest and embedded icon resource used by the build.

## Project Files

- `main.go`: WebView2 window, persistent profile, dark frame, and notifications.
- `app.manifest`: Windows DPI and application manifest.
- `resource.rc`: Windows icon and manifest resource definitions.
- `icon.ico`: Application icon.

## Privacy

This app loads WhatsApp Web directly. Chat data and authentication state are handled by WhatsApp Web and stored locally in the WebView2 profile above. This project is not affiliated with WhatsApp or Meta.

## License

No license has been declared yet.
