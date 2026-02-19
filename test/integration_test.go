//go:build !wasm

package site_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tinywasm/site"
)

type mockIntegrationHandler struct{}

func (h *mockIntegrationHandler) HandlerName() string { return "integration-module" }
func (h *mockIntegrationHandler) RenderHTML() string  { return "module\n<div>Integration</div>" }
func (h *mockIntegrationHandler) RenderCSS() string   { return ".int { color: blue; }" }
func (h *mockIntegrationHandler) AllowedRoles(action byte) []byte {
	return []byte{'*'}
}

func TestIntegration_Render(t *testing.T) {
	site.TestResetHandler()

	err := site.RegisterHandlers(&mockIntegrationHandler{})
	if err != nil {
		t.Fatalf("RegisterHandlers failed: %v", err)
	}

	mux := http.NewServeMux()
	if err := site.Render(mux); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Verify Routes/Assets were registered
	// Use httptest to verify response

	// Create a request
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK for root, got %d", rr.Code)
	}
}
