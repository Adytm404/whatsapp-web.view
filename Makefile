.PHONY: build build-linux build-windows build-macos test test-linux clean lint all

APP_NAME := whatsapp-webview
GO_FLAGS := -ldflags="-s -w"
VERSION ?= dev

# ──────────────────────────────────────────────
# Build targets
# ──────────────────────────────────────────────

all: lint build test

build: build-linux-amd64

build-linux-amd64:
	@echo "→ Building Linux amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o dist/$(APP_NAME)-linux-amd64 .
	@ls -lh dist/$(APP_NAME)-linux-amd64

build-linux-arm64:
	@echo "→ Building Linux arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GO_FLAGS) -o dist/$(APP_NAME)-linux-arm64 .
	@ls -lh dist/$(APP_NAME)-linux-arm64

build-windows:
	@echo "→ Building Windows amd64..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) -H windowsgui -o dist/$(APP_NAME)-windows-amd64.exe .
	@ls -lh dist/$(APP_NAME)-windows-amd64.exe

build-macos-arm64:
	@echo "→ Building macOS arm64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS) -o dist/$(APP_NAME)-macos-arm64 .
	@ls -lh dist/$(APP_NAME)-macos-arm64

build-macos-amd64:
	@echo "→ Building macOS amd64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS) -o dist/$(APP_NAME)-macos-amd64 .
	@ls -lh dist/$(APP_NAME)-macos-amd64

build-all: build-linux-amd64 build-linux-arm64 build-windows build-macos-arm64 build-macos-amd64
	@echo ""
	@echo "=== All builds complete ==="
	@ls -lh dist/

# ──────────────────────────────────────────────
# Test targets
# ──────────────────────────────────────────────

test: test-linux

test-linux:
	@echo "→ Testing Linux binary (Xvfb)..."
	@mkdir -p dist
	@if [ ! -f dist/$(APP_NAME)-linux-amd64 ]; then \
		echo "Binary not found, building..."; \
		$(MAKE) build-linux-amd64; \
	fi
	@bash tools/ci-validate.sh dist/$(APP_NAME)-linux-amd64
	@echo ""
	@echo "→ Running with Xvfb (10s timeout)..."
	@Xvfb :99 -screen 0 1280x720x24 &>/dev/null & XVFB_PID=$$!; \
	sleep 1; \
	DISPLAY=:99 timeout 10 dist/$(APP_NAME)-linux-amd64 2>&1; \
	EXIT_CODE=$$?; \
	kill $$XVFB_PID 2>/dev/null || true; \
	echo ""; \
	echo "Exit code: $$EXIT_CODE"; \
	if [ "$$EXIT_CODE" -eq 124 ]; then \
		echo "✓ PASS: App ran without crashing"; \
	elif [ "$$EXIT_CODE" -eq 0 ]; then \
		echo "✓ PASS: App exited cleanly"; \
	else \
		echo "✗ FAIL: App crashed (exit $$EXIT_CODE)"; \
		exit 1; \
	fi

test-validate:
	@echo "→ Validating all binaries in dist/..."
	@for f in dist/$(APP_NAME)-*; do \
		echo "  $$f: $$(file $$f | cut -d: -f2)"; \
	done

# ──────────────────────────────────────────────
# Lint
# ──────────────────────────────────────────────

lint:
	@echo "→ Go vet..."
	@go vet ./...
	@echo "✓ Lint passed"

# ──────────────────────────────────────────────
# Package
# ──────────────────────────────────────────────

package-deb:
	@echo "→ Building .deb package..."
	@mkdir -p dist
	@if [ ! -f dist/$(APP_NAME)-linux-amd64 ]; then \
		$(MAKE) build-linux-amd64; \
	fi
	@bash packaging/ubuntu/build-deb.sh dist/$(APP_NAME)-linux-amd64 $(VERSION)

# ──────────────────────────────────────────────
# Clean
# ──────────────────────────────────────────────

clean:
	@echo "→ Cleaning..."
	@rm -rf dist/
	@rm -f $(APP_NAME) $(APP_NAME).exe
	@echo "✓ Clean"
