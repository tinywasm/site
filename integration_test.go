//go:build !wasm

package site_test

import (
	"net/http"
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

func TestIntegration_Mount(t *testing.T) {
	// Reset handled indirectly by RegisterHandlers overwriting or appending?
	// The current implementation appends. We might accumulate if we are not careful.
	// But site.go implementation of RegisterHandlers doesn't clear.
	// However, for this test we can just register and check if it appears.

	err := site.RegisterHandlers(&mockIntegrationHandler{})
	if err != nil {
		t.Fatalf("RegisterHandlers failed: %v", err)
	}

	mux := http.NewServeMux()
	if err := site.Mount(mux); err != nil {
		t.Fatalf("Mount failed: %v", err)
	}

	// Verify Routes/Assets were registered
	// We can check if calling the mux returns 404 for expected assets

	req, _ := http.NewRequest("GET", "/style.css", nil)
	// We need a ResponseWriter recorder?
	// Since we are in `site_test` package (external), we can't easily use httptest unless we import it.
	// Assume we can import net/http/httptest
}
