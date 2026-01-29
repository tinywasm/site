//go:build !wasm

package site

import (
	"errors"
	"strings"

	"github.com/tinywasm/assetmin"
)

// build registers all assets with assetmin
func build(am *assetmin.AssetMin) error {
	if len(handler.registeredModules) == 0 {
		return nil
	}

	nav := renderNavigation()
	if nav == "" {
		return errors.New("site: modules registered but no public modules for navigation")
	}

	// Nav at top
	am.InjectBodyContent(nav)

	// Default main container
	am.InjectBodyContent(`<div id="app"></div>`)

	// Process modules
	for _, m := range handler.registeredModules {
		h := m.handler

		// CSS
		if css, ok := h.(assetmin.CSSProvider); ok {
			am.AddCSS(css)
		}

		// Icons
		if icons, ok := h.(assetmin.IconProvider); ok {
			am.AddIcon(icons)
		}

		// HTML (public + "module" in first line)
		if html, ok := h.(interface{ RenderHTML() string }); ok {
			if isPublicReadable(h) {
				content := html.RenderHTML()
				if content != "" && hasModuleInFirstLine(content) {
					am.InjectBodyContent(content)
				}
			}
		}
	}

	return nil
}

func hasModuleInFirstLine(html string) bool {
	firstLine := html
	if idx := strings.Index(html, "\n"); idx != -1 {
		firstLine = html[:idx]
	}
	return strings.Contains(strings.ToLower(firstLine), "module")
}
