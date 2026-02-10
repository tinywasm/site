//go:build !wasm

package site

import (
	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/dom"
)

// trackedComponentsProvider allows site to collect nested components from module builders
// without a direct dependency on the module package.
type trackedComponentsProvider interface {
	TrackedComponents() []dom.HTMLRenderer
}

type titleProvider interface {
	Title() string
}

type accessLevel interface {
	AllowedRoles(action byte) []byte
}

// ssrBuild registers all assets with assetmin
func ssrBuild(am *assetmin.AssetMin) error {
	// 1. Module Discovery: Track components used by registered modules
	for _, m := range handler.registeredModules {
		// If the handler itself is a component, register it
		if comp, ok := m.handler.(dom.HTMLRenderer); ok {
			ssr.componentRegistry.register(comp)
		}

		// If it's a component, trigger its RenderHTML to collect nested components
		// (e.g. if it uses a builder internally)
		if html, ok := m.handler.(dom.HTMLRenderer); ok {
			_ = html.RenderHTML()
		}

		// Now collect everything tracked if the handler provides them
		if tcp, ok := m.handler.(trackedComponentsProvider); ok {
			for _, c := range tcp.TrackedComponents() {
				ssr.componentRegistry.register(c)
			}
		}
	}

	// 2. Asset Injection

	// Inject all collected CSS
	if css := ssr.componentRegistry.collectCSS(); css != "" {
		am.InjectHTML("<style>\n" + css + "</style>\n")
	}

	// Inject all collected JS
	if js := ssr.componentRegistry.collectJS(); js != "" {
		am.InjectHTML("<script>\n" + js + "</script>\n")
	}

	// Inject all collected Icons (Global Sprite)
	for id, svg := range ssr.componentRegistry.collectIcons() {
		am.InjectSpriteIcon(id, svg)
	}

	// 3. Inject Module HTML (public content)
	for _, m := range handler.registeredModules {
		h := m.handler
		if html, ok := h.(dom.HTMLRenderer); ok {
			public := isPublicReadable(h)
			if public {
				content := html.RenderHTML()
				if content != "" {
					am.InjectHTML(content)
				}
			}
		}
	}

	return nil
}

func isPublicReadable(handler any) bool {
	if al, ok := handler.(accessLevel); ok {
		for _, r := range al.AllowedRoles('r') {
			if r == '*' {
				return true
			}
		}
	}
	return false
}
