package site_test

import (
	"reflect"
	"testing"
	"github.com/tinywasm/site"
)

func TestParseRoute(t *testing.T) {
	tests := []struct {
		hash   string
		module string
		params []string
	}{
		{"#users", "users", nil},
		{"#users/123", "users", []string{"123"}},
		{"#users/123/edit", "users", []string{"123", "edit"}},
		{"", "home", nil},
		{"#", "home", nil},
		{"#/users", "users", nil}, // Leading slash
	}

	for _, tt := range tests {
		mod, params := site.TestParseRoute(tt.hash)
		if mod != tt.module {
			t.Errorf("parseRoute(%q) module = %v, want %v", tt.hash, mod, tt.module)
		}

		// Handle nil vs empty slice if necessary, but reflect.DeepEqual handles both as different
		// parseRoute returns nil for params if no params.
		// tt.params is nil or slice.

		if len(params) == 0 && len(tt.params) == 0 {
		    continue
		}

		if !reflect.DeepEqual(params, tt.params) {
			t.Errorf("parseRoute(%q) params = %v, want %v", tt.hash, params, tt.params)
		}
	}
}
