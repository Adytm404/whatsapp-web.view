package platform

// +build windows linux darwin freebsd

type WebViewConfig struct {
	Title    string
	Width    int
	Height   int
	DataPath string
	Debug    bool
}

type ResizeHint int

const (
	HintNone ResizeHint = 0
	HintMin  ResizeHint = 1
	HintMax  ResizeHint = 2
)

type WebView interface {
	SetTitle(title string)
	SetSize(width, height int, hint ResizeHint)
	Navigate(url string)
	Init(script string)
	Bind(name string, fn interface{}) error
	Run()
	Destroy()
	Window() uintptr
}
