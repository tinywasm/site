//go:build !wasm

package main

import (
	"flag"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/site"
	"github.com/tinywasm/site/example/modules"
)

func main() {
	port := flag.String("port", "6060", "server port")
	_ = flag.String("public-dir", "./public", "public directory")
	flag.Parse()

	// 1. Register Handlers
	if err := site.RegisterHandlers(modules.Init()...); err != nil {
		fmt.Println("Error registering handlers:", err)
		return
	}

	// 2. Serve Site (One-liner)
	fmt.Println("Server running http://localhost:" + *port)
	if err := site.Serve(":" + *port); err != nil {
		fmt.Println("Error serving site:", err)
	}
}
