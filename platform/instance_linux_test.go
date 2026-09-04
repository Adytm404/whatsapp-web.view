package platform

import (
	"os"
	"path/filepath"
	"testing"
)

// TestInstanceLockSingleAcquire tests that a single lock can be acquired.
func TestInstanceLockSingleAcquire(t *testing.T) {
	// Use temp dir to avoid polluting user config
	origConfig := os.Getenv("XDG_CONFIG_HOME")
	tempDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Setenv("XDG_CONFIG_HOME", origConfig)

	lock := NewInstanceLock("testlock")
	if !lock.Acquire() {
		t.Fatal("Failed to acquire lock on first attempt")
	}
	defer lock.Release()

	// Lock file should exist
	lockPath := filepath.Join(tempDir, "whatsapp-webview", "lock")
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatalf("Lock file not created: %v", err)
	}
}

// TestInstanceLockSecondBlocks tests that a second instance is blocked.
func TestInstanceLockSecondBlocks(t *testing.T) {
	origConfig := os.Getenv("XDG_CONFIG_HOME")
	tempDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Setenv("XDG_CONFIG_HOME", origConfig)

	lock1 := NewInstanceLock("testlock")
	lock2 := NewInstanceLock("testlock")

	if !lock1.Acquire() {
		t.Fatal("First lock failed to acquire")
	}
	defer lock1.Release()

	if lock2.Acquire() {
		t.Fatal("Second lock should have failed to acquire while first is held")
	}

	lock1.Release()

	// After release, second should be able to acquire
	if !lock2.Acquire() {
		t.Fatal("Second lock failed to acquire after first was released")
	}
	lock2.Release()
}

// TestInstanceLockRelease tests that the lock is properly released.
func TestInstanceLockRelease(t *testing.T) {
	origConfig := os.Getenv("XDG_CONFIG_HOME")
	tempDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Setenv("XDG_CONFIG_HOME", origConfig)

	lock := NewInstanceLock("testlock")
	if !lock.Acquire() {
		t.Fatal("Failed to acquire lock")
	}

	lock.Release()

	// Should be able to re-acquire
	if !lock.Acquire() {
		t.Fatal("Failed to re-acquire after release")
	}
	lock.Release()

	// Lock file should be removed after release
	lockPath := filepath.Join(tempDir, "whatsapp-webview", "lock")
	if _, err := os.Stat(lockPath); err == nil {
		t.Fatal("Lock file should be removed after release")
	}
}
