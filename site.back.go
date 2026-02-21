//go:build !wasm

package site

import (
	"os"
	"reflect"
	"strings"

	"github.com/tinywasm/dom"
)

// ssrState holds SSR-specific state
type ssrState struct {
	assetRegister     assetRegister
	componentRegistry *ssrComponentRegistry
}

var ssr = &ssrState{
	componentRegistry: &ssrComponentRegistry{},
}

func init() {
	env := os.Getenv("APP_ENV")
	if env == "development" || env == "dev" {
		config.DevMode = true
		handler.DevMode = true
	}
	for _, arg := range os.Args {
		if arg == "-dev" {
			config.DevMode = true
			handler.DevMode = true
			break
		}
	}
}

type ssrComponentRegistry struct {
	// registered tracks components by type to avoid duplicate asset collection
	registered map[reflect.Type]dom.Component
}

func (r *ssrComponentRegistry) register(c dom.Component) {
	if r.registered == nil {
		r.registered = make(map[reflect.Type]dom.Component)
	}
	if c == nil {
		return
	}
	t := reflect.TypeOf(c)
	if _, exists := r.registered[t]; !exists {
		r.registered[t] = c
	}
}

// collectCSS generates a single CSS string from all registered components.
func (r *ssrComponentRegistry) collectCSS() string {
	var sb strings.Builder
	for _, c := range r.registered {
		if prov, ok := c.(dom.CSSProvider); ok {
			css := prov.RenderCSS()
			if css != "" {
				sb.WriteString(css)
				sb.WriteString("\n")
			}
		}
	}
	return sb.String()
}

// collectIcons extracts all icons from registered components.
func (r *ssrComponentRegistry) collectIcons() map[string]string {
	icons := make(map[string]string)
	for _, c := range r.registered {
		if prov, ok := c.(dom.IconSvgProvider); ok {
			for id, svg := range prov.IconSvg() {
				icons[id] = svg
			}
		}
	}
	return icons
}

// collectJS generates a single JS string from all registered components.
func (r *ssrComponentRegistry) collectJS() string {
	var sb strings.Builder
	for _, c := range r.registered {
		if prov, ok := c.(dom.JSProvider); ok {
			js := prov.RenderJS()
			if js != "" {
				sb.WriteString(js)
				sb.WriteString("\n")
			}
		}
	}
	return sb.String()
}
