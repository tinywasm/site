//go:build !wasm

package site

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

// renderNavigation generates nav HTML from registered public modules
func renderNavigation() string {
	var html fmt.Conv
	count := 0
	for _, m := range handler.registeredModules {
		if !isPublicReadable(m.handler) || m.name == "help" {
			continue
		}
		displayName := m.name
		if tp, ok := m.handler.(interface{ Title() string }); ok {
			displayName = tp.Title()
		}

		iconID := ""
		if prov, ok := m.handler.(dom.IconSvgProvider); ok {
			for id := range prov.IconSvg() {
				iconID = id
				break
			}
		}

		if count == 0 {
			html.Write("<nav><ul>\n")
		}
		count++

		html.Write("<li><a href='#")
		html.Write(m.name)
		html.Write("' id='nav-")
		html.Write(m.name)
		html.Write("'>")
		if iconID != "" {
			html.Write("<svg><use href='#")
			html.Write(iconID)
			html.Write("'></use></svg> ")
		}
		html.Write(displayName)
		html.Write("</a></li>\n")
	}

	if count > 0 {
		html.Write("</ul></nav>\n")
	}

	return html.String()
}

func isPublicReadable(handler any) bool {
	if al, ok := handler.(dom.AccessLevel); ok {
		for _, r := range al.AllowedRoles('r') {
			if r == '*' {
				return true
			}
		}
	}
	return false
}
