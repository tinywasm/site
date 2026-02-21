//go:build !wasm

package main

import (
	"os"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/site"
	"github.com/tinywasm/site/example/modules"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Production: configure DB, user identity, and roles (uncomment and adapt):
	// site.SetDB(&db.Adapter{DB: openDB()})
	// site.SetUserID(func(data ...any) string {
	// 	for _, d := range data {
	// 		if req, ok := d.(*http.Request); ok {
	// 			return req.Header.Get("X-User-ID")
	// 		}
	// 	}
	// 	return ""
	// })
	// site.CreateRole('a', "Admin",   "Full system access")
	// site.CreateRole('e', "Editor",  "Content management")
	// site.CreateRole('v', "Visitor", "Read-only access")

	if err := site.RegisterHandlers(modules.Init()...); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening on :" + port)
	if err := site.Serve(":" + port); err != nil {
		fmt.Println(err)
	}
}
