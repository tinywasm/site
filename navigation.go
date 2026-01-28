//go:build !wasm

package site

import (
	"github.com/tinywasm/fmt"
)

// renderNavigation generates nav HTML from registered public modules
func renderNavigation() string {
	var links []string
	for _, m := range handler.registeredModules {
		if !isPublicReadable(m.handler) {
			continue
		}
		displayName := m.name
		if dn, ok := m.handler.(interface{ DisplayName() string }); ok {
			displayName = dn.DisplayName()
		}
		link := fmt.Sprintf(`<a href="/%s">%s</a>`, m.name, displayName)
		links = append(links, link)
	}
	if len(links) == 0 {
		return ""
	}
	return fmt.Sprintf(`<nav class="module-nav">%s</nav>`, fmt.Convert(links).Join("").String())
}

func isPublicReadable(handler any) bool {
	if al, ok := handler.(interface{ AllowedRoles(byte) []byte }); ok {
		for _, r := range al.AllowedRoles('r') {
			if r == '*' {
				return true
			}
		}
	}
	return false
}
