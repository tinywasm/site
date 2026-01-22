//go:build !wasm

package site

import (
	"strings"
	"testing"

	"github.com/tinywasm/assetmin"
)

type mockHandler struct {
	name string
	html string
	css  string
	role byte
}

func (h *mockHandler) HandlerName() string { return h.name }
func (h *mockHandler) RenderHTML() string  { return h.html }
func (h *mockHandler) RenderCSS() string   { return h.css }
func (h *mockHandler) AllowedRoles(action byte) []byte {
	if action == 'r' && h.role == '*' {
		return []byte{'*'}
	}
	return []byte{'u'}
}

func TestSite_RegistrationFlow(t *testing.T) {
	// Reset global state for test
	registeredModules = nil

	h1 := &mockHandler{
		name: "nav-module",
		html: "module\n<nav>Nav</nav>",
		css:  ".nav { color: red; }",
		role: '*',
	}
	h2 := &mockHandler{
		name: "private-module",
		html: "module\n<div>Private</div>",
		role: 'u',
	}
	h3 := &mockHandler{
		name: "no-marker-module",
		html: "<div>No marker</div>",
		role: '*',
	}

	RegisterHandlers(h1, h2, h3)

	if len(GetModules()) != 3 {
		t.Errorf("expected 3 modules, got %d", len(GetModules()))
	}

	// Test renderNavigation
	nav := RenderNavigation()
	if !strings.Contains(nav, "nav-module") {
		t.Error("navigation should contain nav-module")
	}
	if strings.Contains(nav, "private-module") {
		t.Error("navigation should NOT contain private-module")
	}

	// Test Build
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: "./test_public",
	})

	err := Build(am)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Internal state check for assetmin (needs bypass or checking exported state)
	// Since we can't easily check private fields of assetmin from another package,
	// we assume that if no errors were returned, the calls happened.
}

func TestSite_Validation(t *testing.T) {
	registeredModules = nil

	// Register only private module
	h := &mockHandler{name: "private", role: 'u'}
	RegisterHandlers(h)

	am := assetmin.NewAssetMin(&assetmin.Config{})
	err := Build(am)
	if err == nil {
		t.Error("Build should fail if no public modules for navigation")
	}
	if !strings.Contains(err.Error(), "no public modules for navigation") {
		t.Errorf("unexpected error: %v", err)
	}
}
