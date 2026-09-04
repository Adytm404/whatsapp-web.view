#!/bin/bash
set -e

VERSION="${1:-1.0.0}"
APP_NAME="whatsapp-webview"

echo "=== Building WhatsApp Webview v${VERSION} ==="

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Go version: ${GO_VERSION}"

# Build Linux (primary target)
echo ""
echo "[1/4] Building Linux amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o "dist/${APP_NAME}-linux-amd64" .
echo "    -> dist/${APP_NAME}-linux-amd64 ($(stat -c%s "dist/${APP_NAME}-linux-amd64" | numfmt --to=iec))"

# Build Linux ARM64
echo "[2/4] Building Linux arm64..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o "dist/${APP_NAME}-linux-arm64" .
echo "    -> dist/${APP_NAME}-linux-arm64 ($(stat -c%s "dist/${APP_NAME}-linux-arm64" | numfmt --to=iec))"

# Build Windows
echo "[3/4] Building Windows amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
    -ldflags="-H windowsgui -s -w -X main.version=${VERSION}" \
    -o "dist/${APP_NAME}-windows-amd64.exe" .
echo "    -> dist/${APP_NAME}-windows-amd64.exe ($(stat -c%s "dist/${APP_NAME}-windows-amd64.exe" | numfmt --to=iec))"

# Build macOS (Apple Silicon)
echo "[4/4] Building macOS arm64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o "dist/${APP_NAME}-macos-arm64" .
echo "    -> dist/${APP_NAME}-macos-arm64 ($(stat -c%s "dist/${APP_NAME}-macos-arm64" | numfmt --to=iec))"

echo ""
echo "=== Build complete ==="
echo "All binaries in dist/"
ls -lh dist/
