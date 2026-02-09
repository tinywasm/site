//go:build !wasm

package site

import (
	"strings"
)

// cssRoot holds the CSS variables for the root element.
type cssRoot struct {
	palette ColorPalette
}

// SetColorPalette sets the global color palette for the site.
func SetColorPalette(p ColorPalette) {
	ssr.cssRoot.palette = p
}

// renderPalette generates the CSS string for the root variables.
func (c *cssRoot) renderPalette() string {
	p := c.palette

	if p.Primary == "" {
		p = ColorPalette{
			Primary:    "#ffffff",
			Secondary:  "#7c3aed",
			Tertiary:   "#94a3b8",
			Quaternary: "#1e293b",
			Gray:       "#f8fafc",
			Selection:  "#a78bfa",
			Hover:      "#6d28d9",
			Success:    "#10b981",
			Error:      "#ef4444",
		}
	}

	var sb strings.Builder
	sb.WriteString(":root {\n")
	sb.WriteString("    /* Colors */\n")

	writeVar(&sb, "--color-primary", p.Primary)
	writeVar(&sb, "--color-secondary", p.Secondary)
	writeVar(&sb, "--color-tertiary", p.Tertiary)
	writeVar(&sb, "--color-quaternary", p.Quaternary)
	writeVar(&sb, "--color-gray", p.Gray)
	writeVar(&sb, "--color-selection", p.Selection)
	writeVar(&sb, "--color-hover", p.Hover)
	writeVar(&sb, "--color-success", p.Success)
	writeVar(&sb, "--color-error", p.Error)

	sb.WriteString("}\n")
	return sb.String()
}

func writeVar(sb *strings.Builder, name, value string) {
	if value != "" {
		sb.WriteString("    " + name + ": " + value + ";\n")
	}
}
