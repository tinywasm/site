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

	// 2. Mount Active Component
	if err := site.Mount("root"); err != nil {
		fmt.Println("Error mounting site:", err)
		return
	}

	fmt.Println("Site mounted successfully ok")

	// 3. Connect responses to handlers (via crudp)
	site.GetCrudP().InitClient()

	select {}
}
