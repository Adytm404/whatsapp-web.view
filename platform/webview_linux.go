//go:build linux

package platform

import (
	"log"

	"github.com/crgimenes/glaze"
	_ "github.com/crgimenes/glaze/embedded"
)

type linuxWebView struct {
	w glaze.WebView
}

func NewWebView(cfg WebViewConfig) WebView {
	w, err := glaze.New(cfg.Debug)
	if err != nil {
		log.Fatalf("Gagal inisialisasi WebView: %v", err)
	}
	w.SetTitle(cfg.Title)
	w.SetSize(cfg.Width, cfg.Height, glaze.HintNone)
	return &linuxWebView{w: w}
}

func (v *linuxWebView) SetTitle(title string) {
	v.w.SetTitle(title)
}

func (v *linuxWebView) SetSize(width, height int, hint ResizeHint) {
	v.w.SetSize(width, height, glaze.HintNone)
}

func (v *linuxWebView) Navigate(url string) {
	v.w.Navigate(url)
}

func (v *linuxWebView) Init(script string) {
	v.w.Init(script)
}

func (v *linuxWebView) Bind(name string, fn interface{}) error {
	return v.w.Bind(name, fn)
}

func (v *linuxWebView) Run() {
	v.w.Run()
}

func (v *linuxWebView) Destroy() {
	v.w.Destroy()
}

func (v *linuxWebView) Window() uintptr {
	return uintptr(v.w.Window())
}
