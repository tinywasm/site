package site

import (
	"github.com/tinywasm/dom"
)

// Module represents a top-level site module.
// It combines a DOM component with site-specific identifiers.
type Module interface {
	dom.Component
	HandlerName() string
	ModuleTitle() string
}

// Parameterized modules can receive route parameters
type Parameterized interface {
	SetParams(params []string)
}

// ModuleLifecycle provides hooks for navigation events
type ModuleLifecycle interface {
	BeforeNavigateAway() bool // Return false to cancel navigation
	AfterNavigateTo()         // Called after module is mounted
}

// ColorPalette defines the global colors for the site.
type ColorPalette struct {
	Primary    string
	Secondary  string
	Tertiary   string
	Quaternary string
	Gray       string
	Selection  string
	Hover      string
	Success    string
	Error      string
}
