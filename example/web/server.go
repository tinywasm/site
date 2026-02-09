//go:build !wasm

package main

import (
	"flag"
	"net/http"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/site"
	"github.com/tinywasm/site/example/modules"
)

func main() {
	port := flag.String("port", "6060", "server port")
	_ = flag.String("public-dir", "./public", "public directory")
	flag.Parse()

	mux := http.NewServeMux()
	if err := site.RegisterHandlers(modules.Init()...); err != nil {
		fmt.Println("Error registering handlers:", err)
		return
	}

	// 4. Render Site (Register routes and build assets)
	if err := site.Render(mux); err != nil {
		fmt.Println("Error rendering site:", err)
		return
	}

	fmt.Println("Server running http://localhost:" + *port)
	http.ListenAndServe(":"+*port, mux)
}
