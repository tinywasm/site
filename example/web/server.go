//go:build !wasm

package main

import (
	"flag"
	"net/http"

	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/site"
	"github.com/tinywasm/site/example/modules"
)

func main() {
	port := flag.String("port", "6060", "server port")
	_ = flag.String("public-dir", "./public", "public directory")
	flag.Parse()

	mux := http.NewServeMux()

	// 1. Register Handlers (orquestador)
	if err := site.RegisterHandlers(modules.Init()...); err != nil {
		fmt.Println("Error registering handlers:", err)
		return
	}

	// 2. Configure AssetMin
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: "./public",
	})

	// 3. Orchestrate Build
	if err := site.Build(am); err != nil {
		fmt.Println("Error building site:", err)
		return
	}

	// 4. Register Routes
	site.GetCrudP().RegisterRoutes(mux)

	fmt.Println("Server running at http://localhost:" + *port)
	http.ListenAndServe(":"+*port, mux)
}
