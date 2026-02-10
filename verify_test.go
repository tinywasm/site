package site

import (
	"testing"
)

func TestConfiguration(t *testing.T) {
	SetCacheSize(5)
	if config.CacheSize != 5 {
		t.Errorf("Expected CacheSize 5, got %d", config.CacheSize)
	}

	SetDefaultRoute("dashboard")
	if config.DefaultRoute != "dashboard" {
		t.Errorf("Expected DefaultRoute dashboard, got %s", config.DefaultRoute)
	}

	SetOutputDir("./dist")
	if config.OutputDir != "./dist" {
		t.Errorf("Expected OutputDir ./dist, got %s", config.OutputDir)
	}

	SetDevMode(true)
	if !config.DevMode {
		t.Errorf("Expected DevMode true, got %v", config.DevMode)
	}
}

func TestParseRouteLogic(t *testing.T) {
	SetDefaultRoute("home")

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
		mod, params := parseRoute(tt.hash)
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
