# Contributing to WhatsApp Webview

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites
- Go 1.22 or newer
- Git

### Linux Development
```bash
# Install dependencies (Ubuntu/Debian)
sudo apt install libwebkit2gtk-4.1-0 libgtk-3-0

# Install dependencies (Arch Linux)
sudo pacman -S webkit2gtk-4.1 gtk3

# Clone the repository
git clone git@github.com:sandikodev/whatsapp-webview.git
cd whatsapp-webview

# Build
CGO_ENABLED=0 go build -o whatsapp-webview .

# Run
./whatsapp-webview
```

## Code Structure

- `main.go` - Shared entry point (platform-agnostic)
- `platform/` - Platform-specific implementations
  - `*_linux.go` - Linux implementations
  - `*_windows.go` - Windows implementations
  - `*_darwin.go` - macOS implementations (planned)

## Adding New Platforms

1. Create new platform files in `platform/` with appropriate build tags
2. Implement required interfaces: `WebView`, `Notifier`, `InstanceLock`
3. Add User-Agent constant in `ua_<platform>.go`
4. Test on the target platform
5. Submit a pull request

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Commit Convention

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
