//go:build wasm

package site_test

import (
	"testing"

	"github.com/tinywasm/dom"
	"github.com/tinywasm/site"
)

type TestModule struct {
	name         string
	beforeCalled bool
	afterCalled  bool
	params       []string
}

func (m *TestModule) HandlerName() string { return m.name }
func (m *TestModule) ModuleTitle() string { return "Test" }

// Implement dom.Component interface manually
func (m *TestModule) RenderHTML() string        { return "" }
func (m *TestModule) RenderCSS() string         { return "" }
func (m *TestModule) GetID() string             { return m.name }
func (m *TestModule) SetID(id string)           { m.name = id }
func (m *TestModule) Children() []dom.Component { return nil }

func (m *TestModule) BeforeNavigateAway() bool {
	m.beforeCalled = true
	return true
}

func (m *TestModule) AfterNavigateTo() {
	m.afterCalled = true
}

func (m *TestModule) SetParams(params []string) {
	m.params = params
}

func TestModuleLifecycle(t *testing.T) {
	// Setup
	m1 := &TestModule{name: "m1"}
	m2 := &TestModule{name: "m2"}

	// Reset global state
	site.TestResetHandler()
	site.TestResetWasm()

	site.RegisterHandlers(m1, m2)

	// Start with m1
	err := site.Navigate("app", "#m1")
	if err != nil {
		t.Logf("Navigate failed: %v", err)
		return
	}

	if !m1.afterCalled {
		t.Error("m1.AfterNavigateTo should be called")
	}

	// Navigation to m2
	err = site.Navigate("app", "#m2")
	if err != nil {
		t.Logf("Navigate failed: %v", err)
		return
	}

	if !m1.beforeCalled {
		t.Error("m1.BeforeNavigateAway should be called")
	}
	if !m2.afterCalled {
		t.Error("m2.AfterNavigateTo should be called")
	}

	// Navigation with params
	err = site.Navigate("app", "#m2/123")
	if err != nil {
		t.Logf("Navigate failed: %v", err)
		return
	}

	if len(m2.params) != 1 || m2.params[0] != "123" {
		t.Errorf("m2 params = %v, want [123]", m2.params)
	}
}
