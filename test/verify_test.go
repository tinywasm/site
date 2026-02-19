package site_test

import (
	"testing"
	"github.com/tinywasm/site"
)

func TestConfiguration(t *testing.T) {
	site.SetCacheSize(5)
	if site.TestGetConfig().CacheSize != 5 {
		t.Errorf("Expected CacheSize 5, got %d", site.TestGetConfig().CacheSize)
	}

	site.SetDefaultRoute("dashboard")
	if site.TestGetConfig().DefaultRoute != "dashboard" {
		t.Errorf("Expected DefaultRoute dashboard, got %s", site.TestGetConfig().DefaultRoute)
	}

	site.SetOutputDir("./dist")
	if site.TestGetConfig().OutputDir != "./dist" {
		t.Errorf("Expected OutputDir ./dist, got %s", site.TestGetConfig().OutputDir)
	}

	site.SetDevMode(true)
	if !site.TestGetConfig().DevMode {
		t.Errorf("Expected DevMode true, got %v", site.TestGetConfig().DevMode)
	}
}

func TestParseRouteLogic(t *testing.T) {
	site.SetDefaultRoute("home")

	tests := []struct {
		hash   string
		module string
		params []string
	}{
		{"#users", "users", []string{}},
		{"#users/123", "users", []string{"123"}},
		{"#users/123/edit", "users", []string{"123", "edit"}},
		{"", "home", []string{}},
		{"#", "home", []string{}},
		{"#/users", "users", []string{}},
	}

	for _, tt := range tests {
		mod, params := site.TestParseRoute(tt.hash)
		if mod != tt.module {
			t.Errorf("parseRoute(%q) module = %v, want %v", tt.hash, mod, tt.module)
		}
		if len(params) != len(tt.params) {
			t.Errorf("parseRoute(%q) params len = %d, want %d", tt.hash, len(params), len(tt.params))
		}
		for i, p := range params {
			if p != tt.params[i] {
				t.Errorf("parseRoute(%q) param[%d] = %v, want %v", tt.hash, i, p, tt.params[i])
			}
		}
	}
}
