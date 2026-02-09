package site

import "github.com/tinywasm/dom"

type mockHandler struct {
	name string
	html string
	css  string
	role byte
}

func (h *mockHandler) HandlerName() string       { return h.name }
func (h *mockHandler) ModuleTitle() string       { return h.name }
func (h *mockHandler) RenderHTML() string        { return h.html }
func (h *mockHandler) RenderCSS() string         { return h.css }
func (h *mockHandler) ID() string                { return h.name }
func (h *mockHandler) SetID(id string)           {}
func (h *mockHandler) Children() []dom.Component { return nil }
func (h *mockHandler) AllowedRoles(action byte) []byte {
	if action == 'r' && h.role == '*' {
		return []byte{'*'}
	}
	return []byte{'u'}
}
