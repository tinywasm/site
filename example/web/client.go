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

	// 2. Mount Site (Blocks automatically)
	fmt.Println("Mounting site...")
	if err := site.Mount("app"); err != nil {
		fmt.Println("Error mounting site:", err)
	}
	// No need for select{} anymore
}
