# Contributing to WhatsApp Webview

Thank you for your interest in contributing! This guide will help you get started.

## CI/CD Pipeline

Before opening a PR, your changes are automatically verified by **GitHub Actions**:

1. **lint** — `go vet`, `go mod verify`, `go mod tidy` integrity
2. **build** — Matrix build for all 5 supported platforms
3. **test-linux** — Real WebView smoke test under **Xvfb** (virtual display)
4. **validate** — Binary format/size checks for all artifacts
5. **package** — `.deb` package built from the tested binary

> ⚠️ **Your PR will not merge if any CI job fails.** All stages are mandatory.

### Run CI checks locally first

```bash
# Full CI simulation (lint + build + test)
make all

# Or step by step:
make lint          # Stage 1: lint
make build-all     # Stage 2: build all platforms
make test          # Stage 3-4: test + validate Linux binary
```

### CI badges

```markdown
[![CI](https://github.com/sandikodev/whatsapp-webview/actions/workflows/ci.yml/badge.svg)](https://github.com/sandikodev/whatsapp-webview/actions/workflows/ci.yml)
```

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

# Build (all platforms + validate)
./build.sh 1.0.0

# Or just Linux:
make build

# Lint
make lint

# Test (needs Xvfb + WebKitGTK for full smoke test)
make test

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

1. Create new platform files in `platform/` with appropriate build tags (`//go:build <os>`)
2. Implement required interfaces: `WebView`, `Notifier`, `InstanceLock`
3. Add User-Agent constant in `ua_<platform>.go`
4. Add the platform to the build matrix in `.github/workflows/ci.yml`
5. Add packaging for the target package manager
6. Test on the target platform (or extend CI if a runner exists)
7. Submit a pull request

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Run CI checks locally: `make all`
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request on GitHub
7. Ensure all CI checks pass (lint, build matrix, Xvfb test, validation, package)

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
