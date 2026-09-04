#!/bin/bash
set -euo pipefail

VERSION="${1:-dev}"
APP_NAME="whatsapp-webview"
BUILD_DIR="dist"
LDFLAGS="-s -w -X main.version=${VERSION}"

info()  { echo -e "\033[1;36m[INFO]\033[0m $*"; }
ok()    { echo -e "\033[1;32m[ OK ]\033[0m $*"; }
warn()  { echo -e "\033[1;33m[WARN]\033[0m $*"; }
fail()  { echo -e "\033[1;31m[FAIL]\033[0m $*"; exit 1; }

GO_VERSION=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//') || fail "Go not found. Install Go first."
info "Build ${APP_NAME} v${VERSION} (Go ${GO_VERSION})"

mkdir -p "${BUILD_DIR}"

build() {
    local target="$1" goos="$2" goarch="$3" extra_ldflags="${4:-}"
    local artifact="${BUILD_DIR}/${APP_NAME}-${target}"
    local all_ldflags="${LDFLAGS}"

    [ -n "${extra_ldflags}" ] && all_ldflags="${all_ldflags} ${extra_ldflags}"

    info "Building ${target} (${goos}/${goarch})..."
    CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
        go build -trimpath -ldflags="${all_ldflags}" \
        -o "${artifact}" .
    ok "Created ${artifact} ($(du -h "${artifact}" | cut -f1))"
}

# Platform matrix
build "linux-amd64"     linux   amd64
build "linux-arm64"     linux   arm64
build "windows-amd64"   windows amd64 "-H windowsgui"
build "macos-arm64"     darwin  arm64
build "macos-amd64"     darwin  amd64

info "Validating binaries..."
for f in "${BUILD_DIR}"/${APP_NAME}-*; do
    bash tools/ci-validate.sh "${f}" || fail "Validation failed for ${f}"
done

info "Generating checksums..."
cd "${BUILD_DIR}"
if command -v sha256sum >/dev/null 2>&1; then
    sha256sum ${APP_NAME}-* > SHA256SUMS.txt
elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 ${APP_NAME}-* > SHA256SUMS.txt
else
    warn "No sha256sum available, skipping checksums"
fi
cd - >/dev/null

ok "Build complete!"
info "Binaries in ${BUILD_DIR}/:"
ls -lh "${BUILD_DIR}"

if [ -f "${BUILD_DIR}/SHA256SUMS.txt" ]; then
    info "Checksums:"
    cat "${BUILD_DIR}/SHA256SUMS.txt"
fi
