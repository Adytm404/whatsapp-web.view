#!/bin/bash
set -e

# Usage:
#   ./build-deb.sh [VERSION]                    # builds binary, then packages
#   ./build-deb.sh <binary-path> <version>      # packages existing binary
VERSION="${2:-${1:-1.0.0}}"
BINARY_INPUT="${1:-}"  # if set and is a file path, use it instead of building
ARCH="amd64"
PKG_NAME="whatsapp-webview"
BUILD_DIR="build/deb"

echo "=== Building WhatsApp Webview .deb (v${VERSION}) ==="

mkdir -p "${BUILD_DIR}/usr/bin"
mkdir -p "${BUILD_DIR}/usr/share/applications"
mkdir -p "${BUILD_DIR}/usr/share/icons/hicolor/128x128/apps"
mkdir -p "${BUILD_DIR}/DEBIAN"

if [ -n "${BINARY_INPUT}" ] && [ -f "${BINARY_INPUT}" ]; then
    echo "[1/5] Using prebuilt binary: ${BINARY_INPUT}"
    cp "${BINARY_INPUT}" "${BUILD_DIR}/usr/bin/${PKG_NAME}"
elif [ -n "${BINARY_INPUT}" ] && ! command -v "${BINARY_INPUT}" >/dev/null 2>&1; then
    echo "ERROR: Binary not found: ${BINARY_INPUT}"
    echo "Usage: ./build-deb.sh [VERSION] or ./build-deb.sh <binary> <version>"
    exit 1
else
    echo "[1/5] Building Go binary..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-s -w" \
        -o "${BUILD_DIR}/usr/bin/${PKG_NAME}" \
        .
fi

chmod 755 "${BUILD_DIR}/usr/bin/${PKG_NAME}"

echo "[2/5] Installing desktop file..."
cp packaging/ubuntu/usr/share/applications/*.desktop \
    "${BUILD_DIR}/usr/share/applications/"

echo "[3/5] Installing icon..."
if [ -f assets/icon.png ]; then
    cp assets/icon.png \
        "${BUILD_DIR}/usr/share/icons/hicolor/128x128/apps/${PKG_NAME}.png"
fi

echo "[4/5] Copying DEBIAN files..."
cp packaging/ubuntu/DEBIAN/* "${BUILD_DIR}/DEBIAN/"
chmod 755 "${BUILD_DIR}/DEBIAN/postinst"
chmod 755 "${BUILD_DIR}/DEBIAN/prerm"

echo "[5/5] Building .deb..."
dpkg-deb --build "${BUILD_DIR}" "${PKG_NAME}_${VERSION}_${ARCH}.deb"

echo ""
echo "Done! Package: ${PKG_NAME}_${VERSION}_${ARCH}.deb"
echo ""
echo "Install with:"
echo "  sudo dpkg -i ${PKG_NAME}_${VERSION}_${ARCH}.deb"
echo "  sudo apt install -f"
