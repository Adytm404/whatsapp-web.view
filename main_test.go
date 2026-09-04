package main

import (
	"strings"
	"testing"
)

func TestGetInitScript(t *testing.T) {
	ua := "X11; Linux x86_64"
	script := getInitScript(ua)

	// Should contain the user agent
	if !strings.Contains(script, ua) {
		t.Error("Init script does not contain User-Agent")
	}

	// Should contain notification polyfill
	if !strings.Contains(script, "sendNativeNotification") {
		t.Error("Init script does not contain notification handler")
	}

	// Should contain dark mode CSS
	if !strings.Contains(script, "darkCSS") {
		t.Error("Init script does not contain dark mode CSS")
	}

	// Should contain permission granting
	if !strings.Contains(script, "'granted'") {
		t.Error("Init script does not grant notification permission")
	}
}

func TestGetInitScriptChromeVersion(t *testing.T) {
	script := getInitScript("Windows NT 10.0")

	// Should have Chrome 133 user agent
	if !strings.Contains(script, "Chrome/133.0.0.0") {
		t.Error("Init script missing Chrome version")
	}

	// Should reference both userAgent and appVersion
	if strings.Count(script, "Chrome/133.0.0.0") < 2 {
		t.Error("User-Agent should appear in both userAgent and appVersion")
	}
}

func TestGetUserDataDir(t *testing.T) {
	t.Log("getUserDataDir should return %LOCALAPPDATA% or ~/.config based on OS")
	dir := getUserDataDir()
	if !strings.Contains(dir, "WhatsAppDesktopLight") {
		t.Errorf("UserDataDir should contain 'WhatsAppDesktopLight', got: %s", dir)
	}
}
