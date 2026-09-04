#!/bin/bash
# ci-validate.sh — Validate a compiled binary
# Usage: ci-validate.sh <binary-path>
set -e

BINARY="${1:?Usage: ci-validate.sh <binary-path>}"

if [ ! -f "${BINARY}" ]; then
  echo "FAIL: File not found: ${BINARY}"
  exit 1
fi

echo "=== Binary Validation: $(basename ${BINARY}) ==="
echo ""

# 1. File exists and is non-empty
SIZE=$(stat -c%s "${BINARY}" 2>/dev/null || stat -f%z "${BINARY}" 2>/dev/null)
if [ "${SIZE}" -eq 0 ]; then
  echo "FAIL: Binary is empty"
  exit 1
fi
echo "✓ Size: ${SIZE} bytes ($(( SIZE / 1024 / 1024 ))MB)"

# 2. File type check
FILE_TYPE=$(file "${BINARY}")
echo "  Type: ${FILE_TYPE}"

# 3. Platform-specific checks
case "${FILE_TYPE}" in
  *"ELF"*"x86-64"*)
    echo "✓ Format: ELF x86-64 (Linux)"
    # Check for static linking
    if echo "${FILE_TYPE}" | grep -q "statically linked"; then
      echo "✓ Linked: Static"
    elif echo "${FILE_TYPE}" | grep -q "dynamically linked"; then
      echo "⚠ Linked: Dynamic (may need runtime libs)"
    fi
    # Check for stripped symbols
    if echo "${FILE_TYPE}" | grep -q "not stripped"; then
      echo "⚠ Symbols: Not stripped (binary larger)"
    else
      echo "✓ Symbols: Stripped"
    fi
    ;;
  *"ELF"*"aarch64"*)
    echo "✓ Format: ELF aarch64 (Linux ARM64)"
    ;;
  *"PE32"*"Windows"*)
    echo "✓ Format: PE32+ (Windows)"
    ;;
  *"Mach-O"*"arm64"*)
    echo "✓ Format: Mach-O arm64 (macOS)"
    ;;
  *"Mach-O"*"x86_64"*)
    echo "✓ Format: Mach-O x86-64 (macOS)"
    ;;
  *)
    echo "⚠ Unknown format: ${FILE_TYPE}"
    ;;
esac

# 4. Size limit (20MB)
if [ "${SIZE}" -gt 20971520 ]; then
  echo "FAIL: Binary too large (>20MB)"
  exit 1
fi
echo "✓ Size within limits"

# 5. Executable check
if [ -x "${BINARY}" ]; then
  echo "✓ Executable: Yes"
else
  echo "⚠ Executable: No (may need chmod +x)"
fi

echo ""
echo "=== Validation passed ==="
