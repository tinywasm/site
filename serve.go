//go:build !wasm

package site

import "net/http"

// Serve starts the server on the given address (one-liner helper).
// It creates a new ServeMux, mounts the site, and listens on the address.
func Serve(addr string) error {
	mux := http.NewServeMux()
	if err := Mount(mux); err != nil {
		return err
	}
	return http.ListenAndServe(addr, mux)
}
