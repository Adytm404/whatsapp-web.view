package platform

// +build linux

const OSUserAgent = "X11; Linux x86_64"

func SetDarkFrame(window uintptr) {
	// Dark mode on Linux is handled via CSS injection in main.go
	// through w.Init() with the dark mode CSS styles.
	// This function is a no-op because glaze doesn't expose
	// the GTK window handle for direct CSS provider injection.
}
