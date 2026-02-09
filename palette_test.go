//go:build !wasm

package site

import (
	"strings"
	"testing"
)

func TestSetColorPalette(t *testing.T) {
	p := ColorPalette{
		Primary: "red",
	}
	SetColorPalette(p)

	if ssr.cssRoot.palette.Primary != "red" {
		t.Errorf("Expected Primary to be red, got %s", ssr.cssRoot.palette.Primary)
	}
}

func TestRenderPalette(t *testing.T) {
	p := ColorPalette{
		Primary:   "#ffffff",
		Secondary: "#000000",
	}
	SetColorPalette(p)

	css := ssr.cssRoot.renderPalette()
	if !strings.Contains(css, "--color-primary: #ffffff;") {
		t.Errorf("Expected CSS to contain primary color, got %s", css)
	}
	if !strings.Contains(css, "--color-secondary: #000000;") {
		t.Errorf("Expected CSS to contain secondary color, got %s", css)
	}
}
