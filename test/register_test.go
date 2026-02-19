//go:build !wasm

package site_test

import (
	"testing"

	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/site"
)

func TestSite_RegistrationFlow(t *testing.T) {
	// Reset global state for test
	site.TestResetHandler()

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

	err := site.RegisterHandlers(h1, h2, h3)
	if err != nil {
		t.Fatalf("RegisterHandlers failed: %v", err)
	}

	// All 3 handlers implement HTMLProvider, so all get registered
	if len(site.TestGetModules()) != 3 {
		t.Errorf("expected 3 modules, got %d", len(site.TestGetModules()))
	}

	// Test build
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: "./test_public",
	})

	err = site.TestSSRBuild(am)
	if err != nil {
		t.Fatalf("ssrBuild failed: %v", err)
	}

	// Internal state check for assetmin (needs bypass or checking exported state)
	// Since we can't easily check private fields of assetmin from another package,
	// we assume that if no errors were returned, the calls happened.
}
