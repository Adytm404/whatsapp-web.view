# WhatsApp Webview Desktop

Lightweight cross-platform desktop wrapper for [WhatsApp Web](https://web.whatsapp.com), built with Go and system WebView.

## Features

- Native WebView window (no Electron)
- Persistent WhatsApp session across app restarts
- Cookies, LocalStorage, IndexedDB, and service-worker data stored in a dedicated profile
- Dark title bar theme (WhatsApp-style dark header)
- Native notification bridge (OS-level notifications)
- Camera and microphone access for WhatsApp voice and video calls
- Single-instance protection (prevents multiple windows)
- High-DPI display support
- Small native executable (~6MB on Linux, ~5MB on Windows)
- No bundled browser engine (uses system WebView)
- Cross-compile from any OS without C toolchain

## Supported Platforms

| Platform | WebView Engine | Status | Package |
|----------|---------------|--------|---------|
| Windows 10/11 | Microsoft Edge WebView2 | ✅ Supported | `.exe` |
| Ubuntu 22.04+ / Debian 12+ | WebKitGTK 4.1 | ✅ Supported | Manual build / `.deb` |
| Arch Linux / Manjaro | WebKitGTK 4.1 | ✅ Supported | `paru` / `yay` |
| Fedora 40+ | WebKitGTK 4.1 | ✅ Supported | Manual build |
| macOS (Sonoma+) | WKWebView | 🔜 Planned | — |
| FreeBSD 14+ | WebKitGTK | 🔜 Planned | — |

---

## Table of Contents

- [Installation](#installation)
  - [Ubuntu / Debian](#ubuntu--debian)
  - [Arch Linux](#arch-linux)
  - [Fedora](#fedora)
  - [Windows](#windows)
- [Build From Source](#build-from-source)
  - [Prerequisites](#prerequisites)
  - [Build Commands](#build-commands)
  - [Build Flags Explained](#build-flags-explained)
- [Usage](#usage)
  - [First Launch](#first-launch)
  - [Session Data](#session-data)
  - [Troubleshooting](#troubleshooting)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [Privacy](#privacy)
- [Credits](#credits)

---

## Installation

### Ubuntu / Debian

#### Option 1: Build from Source (Recommended)

**Step 1: Install Go 1.22+**

Ubuntu's default Go version is too old. Install the latest Go:

```bash
# Download Go 1.22
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz

# Remove old Go and install new
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc for persistence)
export PATH=$PATH:/usr/local/go/bin

# Verify
go version
# Should show: go version go1.22.5 linux/amd64
```

**Step 2: Install WebKitGTK dependencies**

```bash
sudo apt update
sudo apt install -y libwebkit2gtk-4.1-0 libgtk-3-0
```

**Step 3: Clone and build**

```bash
# Clone the repository
git clone https://github.com/sandikodev/whatsapp-webview.git
cd whatsapp-webview

# Build (optimized, no debug symbols)
CGO_ENABLED=0 go build -ldflags="-s -w" -o whatsapp-webview .

# The binary is now ready
ls -lh whatsapp-webview
# -rwxr-xr-x 1 user user 6.1M Sep  4 whatsapp-webview
```

**Step 4: Install system-wide (optional)**

```bash
# Copy binary to /usr/local/bin
sudo cp whatsapp-webview /usr/local/bin/

# Create desktop entry
sudo mkdir -p /usr/share/applications
sudo tee /usr/share/applications/whatsapp-webview.desktop > /dev/null << EOF
[Desktop Entry]
Name=WhatsApp Webview
Comment=Lightweight WhatsApp Desktop Webview
Exec=/usr/local/bin/whatsapp-webview
Icon=whatsapp-webview
Terminal=false
Type=Application
Categories=Network;InstantMessaging;Chat;
Keywords=whatsapp;web;chat;messaging;
StartupWMClass=whatsapp-webview
EOF

# Update desktop database
sudo update-desktop-database /usr/share/applications

# Now you can launch from app menu or terminal
whatsapp-webview
```

**Step 5: Create .deb package (optional)**

```bash
# Build the .deb package
bash build-deb.sh 1.0.0

# Install it
sudo dpkg -i whatsapp-webview_1.0.0_amd64.deb
sudo apt install -f  # Fix any missing dependencies

# Uninstall if needed
sudo dpkg -r whatsapp-webview
```

#### Option 2: Using paru/AUR (when available)

Once the AUR package is published:

```bash
paru -S whatsapp-webview-git
```

#### System Requirements (Ubuntu/Debian)

| Package | Minimum Version | Purpose |
|---------|----------------|---------|
| `libwebkit2gtk-4.1-0` | 2.40+ | WebView engine |
| `libgtk-3-0` | 3.24+ | GTK toolkit |
| `libnotify4` | 0.7+ | Notification bridge (optional) |

Check installed versions:
```bash
dpkg -l | grep -E "libwebkit2gtk-4.1|libgtk-3-0"
```

---

### Arch Linux

#### Option 1: Install via paru (Recommended)

**Step 1: Install paru (if not installed)**

```bash
# Install paru from AUR
sudo pacman -S --needed base-devel
git clone https://aur.archlinux.org/paru.git
cd paru
makepkg -si
```

**Step 2: Install WebKitGTK dependencies**

```bash
sudo pacman -S webkit2gtk-4.1 gtk3
```

**Step 3: Install whatsapp-webview**

```bash
# Install from AUR
paru -S whatsapp-webview-git

# Launch
whatsapp-webview
```

#### Option 2: Build from Source

**Step 1: Install Go**

```bash
sudo pacman -S go
```

**Step 2: Install WebKitGTK**

```bash
sudo pacman -S webkit2gtk-4.1 gtk3
```

**Step 3: Clone and build**

```bash
git clone https://github.com/sandikodev/whatsapp-webview.git
cd whatsapp-webview

# Build
CGO_ENABLED=0 go build -ldflags="-s -w" -o whatsapp-webview .

# Run
./whatsapp-webview
```

**Step 4: Create PKGBUILD (for AUR publishing)**

```bash
# The PKGBUILD is already included in packaging/archlinux/
cp packaging/archlinux/PKGBUILD .

# Build the package
makepkg -si

# This installs whatsapp-webview system-wide
```

#### Option 3: Using .deb (when available)

If a `.deb` package is available, Arch Linux can use it via `debtap`:

```bash
# Install debtap
sudo pacman -S debtap

# Convert .deb to Arch package
debtap whatsapp-webview_1.0.0_amd64.deb

# Install the converted package
sudo pacman -U whatsapp-webview-1.0.0-1-x86_64.pkg.tar.zst
```

> **Note:** Using `.deb` via debtap provides more stable packages as they're tested on Debian-based systems first.

#### System Requirements (Arch Linux)

| Package | Purpose |
|---------|---------|
| `webkit2gtk-4.1` | WebView engine |
| `gtk3` | GTK toolkit |

Check installed:
```bash
pacman -Qi webkit2gtk-4.1 gtk3
```

---

### Fedora

```bash
# Install dependencies
sudo dnf install webkit2gtk4.1 gtk3

# Clone and build
git clone https://github.com/sandikodev/whatsapp-webview.git
cd whatsapp-webview
CGO_ENABLED=0 go build -ldflags="-s -w" -o whatsapp-webview .

# Run
./whatsapp-webview
```

---

### Windows

**Option 1: Download pre-built**

Download `WhatsApp.exe` from the latest [GitHub Release](https://github.com/sandikodev/whatsapp-webview/releases).

**Option 2: Build from source**

```powershell
# Install Go from https://go.dev/dl/

# Clone
git clone https://github.com/sandikodev/whatsapp-webview.git
cd whatsapp-webview

# Build (hides console window)
go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .

# Run
.\WhatsApp.exe
```

**Requirements:**
- Windows 10 or newer
- Microsoft Edge WebView2 Runtime (auto-downloaded if missing)
- WhatsApp account paired with WhatsApp Web

---

## Build From Source

### Prerequisites

| Requirement | Version | Check |
|------------|---------|-------|
| Go | 1.22+ | `go version` |
| Git | Any | `git --version` |
| WebKitGTK (Linux) | 4.1+ | See platform-specific commands |

### Build Commands

```bash
# Linux (x86_64)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o whatsapp-webview .

# Linux (ARM64)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o whatsapp-webview-arm64 .

# Windows (x86_64)
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .

# macOS (Apple Silicon)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o whatsapp-webview-macos .

# macOS (Intel)
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o whatsapp-webview-macos-intel .

# FreeBSD (x86_64)
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o whatsapp-webview-freebsd .
```

### Build Flags Explained

| Flag | Purpose |
|------|---------|
| `CGO_ENABLED=0` | Disable CGo for pure Go build (required for cross-compile) |
| `GOOS=linux` | Target operating system |
| `GOARCH=amd64` | Target architecture |
| `-ldflags="-s -w"` | Strip debug symbols (smaller binary) |
| `-ldflags="-H windowsgui"` | Windows only: hide console window |

### Binary Size Comparison

| Platform | Size | Notes |
|----------|------|-------|
| Linux (stub) | 1.9MB | Without glaze |
| Linux (glaze) | 6.1MB | With WebKitGTK dlopen |
| Windows | 4.8MB | With WebView2 |

---

## Usage

### First Launch

1. Run `whatsapp-webview` (or `WhatsApp.exe` on Windows)
2. A window will open with WhatsApp Web
3. Scan the QR code with your phone
4. Your session is saved automatically
5. Close and reopen — you stay logged in

### Session Data

Profile data is stored at:

| Platform | Path |
|----------|------|
| Windows | `%APPDATA%\WhatsAppDesktopLight\UserData` |
| Linux | `~/.config/WhatsAppDesktopLight/UserData` |
| macOS | `~/Library/Application Support/WhatsAppDesktopLight/UserData` |

> **Important:** Do not delete this folder if you want to keep your login session.

### Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+Q` | Quit application |
| `Alt+F4` | Close window |
| `F11` | Toggle fullscreen (if supported by WebView) |

### Troubleshooting

#### "Failed to initialize WebView" (Linux)

**Cause:** WebKitGTK libraries not installed or wrong version.

**Fix:**
```bash
# Ubuntu/Debian
sudo apt install libwebkit2gtk-4.1-0

# Arch Linux
sudo pacman -S webkit2gtk-4.1

# Verify
ldconfig -p | grep webkit2gtk
```

#### "Segmentation fault" on launch

**Cause:** Conflicting GTK themes or extensions.

**Fix:**
```bash
# Try with default theme
GTK_THEME=Adwaita ./whatsapp-webview

# Or reset GTK settings
rm -rf ~/.config/gtk-3.0
```

#### Notifications not working (Linux)

**Cause:** D-Bus session bus not available.

**Fix:**
```bash
# Check D-Bus
dbus-send --session --print-reply \
  --dest=org.freedesktop.Notifications \
  /org/freedesktop/Notifications \
  org.freedesktop.Notifications.Notify \
  string:"Test" string:"Test notification"

# If fails, install notification daemon
sudo apt install dunst  # or notification-daemon
```

#### WhatsApp shows "unsupported browser"

**Cause:** User-Agent not set correctly.

**Fix:** This shouldn't happen with the built-in UA spoofing. If it does:
```bash
# Check the binary is up to date
git pull
go build -ldflags="-s -w" -o whatsapp-webview .
```

#### Multiple instance protection not working

**Cause:** Lock file from previous crashed instance.

**Fix:**
```bash
# Remove stale lock
rm -f ~/.config/whatsapp-webview/lock

# Or kill any remaining processes
pkill -f whatsapp-webview
```

---

## Project Structure

```
whatsapp-webview/
├── main.go                        # Shared entry point (all platforms)
├── go.mod                         # Go module definition
├── go.sum                         # Dependency checksums
│
├── platform/                      # Platform-specific implementations
│   ├── webview.go                 # WebView interface (shared)
│   ├── webview_linux.go           # Linux: crgimenes/glaze (WebKitGTK)
│   ├── webview_windows.go         # Windows: jchv/go-webview2 (Edge WebView2)
│   ├── notify.go                  # Notifier interface (shared)
│   ├── notify_linux.go            # Linux: D-Bus org.freedesktop.Notifications
│   ├── notify_windows.go          # Windows: go-toast/toast (COM)
│   ├── instance.go                # InstanceLock interface (shared)
│   ├── instance_linux.go          # Linux: flock() file lock
│   ├── instance_windows.go        # Windows: CreateMutexW / FindWindowW
│   ├── frame_linux.go             # Linux: CSS injection dark mode
│   ├── frame_windows.go           # Windows: DwmSetWindowAttribute
│   ├── ua_windows.go              # Windows User-Agent string
│   ├── platform_darwin.go         # macOS: stub (planned)
│   └── platform_freebsd.go        # FreeBSD: stub (planned)
│
├── packaging/                     # Distribution packages
│   ├── ubuntu/
│   │   ├── DEBIAN/
│   │   │   ├── control            # Package metadata
│   │   │   ├── postinst           # Post-install script
│   │   │   └── prerm              # Pre-remove script
│   │   └── usr/share/applications/
│   │       └── whatsapp-webview.desktop
│   └── archlinux/
│       └── PKGBUILD               # AUR build script
│
├── build-deb.sh                   # Build .deb package script
├── CONTRIBUTING.md                 # Contribution guidelines
├── README.md                      # This file
├── .gitignore                     # Git ignore rules
├── app.manifest                   # Windows DPI manifest
├── gen_icon.py                    # Icon generation script
├── icon.ico                       # Windows icon
└── resource.rc                    # Windows resource file
```

---

## Architecture

### Cross-Platform Strategy

The project uses **Go build tags** to compile platform-specific code:

```
main.go (shared)
    │
    ├── platform/webview.go (interface)
    │   ├── webview_linux.go   → crgimenes/glaze → WebKitGTK
    │   └── webview_windows.go → jchv/go-webview2 → Edge WebView2
    │
    ├── platform/notify.go (interface)
    │   ├── notify_linux.go   → godbus/dbus → D-Bus
    │   └── notify_windows.go → go-toast/toast → COM
    │
    └── platform/instance.go (interface)
        ├── instance_linux.go   → syscall.Flock()
        └── instance_windows.go → CreateMutexW
```

### Why crgimenes/glaze for Linux?

| Criteria | glaze | webview/webview | go-webview2 |
|----------|-------|-----------------|-------------|
| CGo-free | ✅ | ❌ | ✅ |
| Linux support | ✅ | ✅ | ❌ |
| macOS support | ✅ | ✅ | ❌ |
| Windows support | ✅ | ✅ | ✅ |
| Cross-compile | ✅ | ❌ | ✅ |
| No bundled libs | ✅ | ❌ | ✅ |

**glaze** uses `dlopen()` to load WebKitGTK at runtime — no compilation against native libraries needed.

### User-Agent Spoofing

The app spoofs the User-Agent to report as Chrome on the target OS:

| Platform | User-Agent |
|----------|------------|
| Windows | `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36...` |
| Linux | `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36...` |
| macOS | `Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/537.36...` |
| FreeBSD | `Mozilla/5.0 (X11; FreeBSD amd64) AppleWebKit/537.36...` |

---

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Start for Contributors

```bash
# Fork the repository on GitHub

# Clone your fork
git clone https://github.com/YOUR_USERNAME/whatsapp-webview.git
cd whatsapp-webview

# Create a feature branch
git checkout -b feat/my-new-feature

# Make your changes, then build and test
CGO_ENABLED=0 go build -ldflags="-s -w" -o whatsapp-webview .

# Commit and push
git add .
git commit -m "feat: add my new feature"
git push origin feat/my-new-feature

# Open a Pull Request on GitHub
```

---

## Privacy

This app loads WhatsApp Web directly. Chat data and authentication state are handled by WhatsApp Web and stored locally in the WebView profile. This project is not affiliated with WhatsApp or Meta.

**Data stored locally:**
- Login session (cookies, tokens)
- Cached messages and media
- User preferences

**Data NOT collected:**
- No analytics
- No telemetry
- No remote logging
- No phone data extraction

---

## Credits

- Original Windows-only project: [Adytm404/whatsapp-web.view](https://github.com/Adytm404/whatsapp-web.view)
- Cross-platform WebView: [crgimenes/glaze](https://github.com/crgimenes/glaze)
- D-Bus notifications: [godbus/dbus](https://github.com/godbus/dbus)
- Windows WebView2: [jchv/go-webview2](https://github.com/jchv/go-webview2)
