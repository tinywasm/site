//go:build wasm

package main

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/site"
	"github.com/tinywasm/site/example/modules"
)

func main() {
	// 1. Register Handlers
	if err := site.RegisterHandlers(modules.Init()...); err != nil {
		fmt.Println("Error registering handlers:", err)
		return
	}

	// 2. Render Site (Hydrate initial, then Render on navigation)
	if err := site.Render(); err != nil {
		fmt.Println("Error rendering site:", err)
		return
	}

	fmt.Println("Site rendered successfully ok")

	select {}
}
