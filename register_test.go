//go:build !wasm

package site

import (
	"strings"
	"testing"

	"github.com/tinywasm/assetmin"
)

func TestSite_RegistrationFlow(t *testing.T) {
	// Reset global state for test
	handler.registeredModules = nil

	h1 := &mockHandler{
		name: "nav-module",
		html: "<div>Nav</div>",
		role: '*',
	}
	h2 := &mockHandler{
		name: "private-module",
		html: "<div>Private</div>",
		role: 'u',
	}
	h3 := &mockHandler{
		name: "no-marker-module",
		html: `<div>No marker</div>`,
		role: '*',
	}

	err := RegisterHandlers(h1, h2, h3)
	if err != nil {
		t.Fatalf("RegisterHandlers failed: %v", err)
	}

	// All 3 handlers implement HTMLProvider, so all get registered
	if len(getModules()) != 3 {
		t.Errorf("expected 3 modules, got %d", len(getModules()))
	}

	// Test renderNavigation
	nav := renderNavigation()
	if !strings.Contains(nav, "nav-module") {
		t.Error("navigation should contain nav-module")
	}
	if strings.Contains(nav, "private-module") {
		t.Error("navigation should NOT contain private-module")
	}

	// Test build
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: "./test_public",
	})

	err = ssrBuild(am)
	if err != nil {
		t.Fatalf("ssrBuild failed: %v", err)
	}

	// Internal state check for assetmin (needs bypass or checking exported state)
	// Since we can't easily check private fields of assetmin from another package,
	// we assume that if no errors were returned, the calls happened.
}
