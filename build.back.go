//go:build !wasm

package site

import (
	"os"

	"github.com/tinywasm/assetmin"
	"github.com/tinywasm/fmt"
)

const staticBuildFlag = "--ssr-static-build"

// BuildStatic renders all registered modules and writes the output
// to outputDir as static HTML/CSS/JS/SVG files using assetmin.
func BuildStatic(outputDir string) error {
	am := assetmin.NewAssetMin(&assetmin.Config{
		OutputDir: outputDir,
	})
	am.EnsureOutputDirectoryExists()
	if err := ssrBuild(am); err != nil {
		return err
	}
	am.SetBuildOnDisk(true)
	return nil
}

// AutoBuild checks os.Args for --ssr-static-build <dir>.
// If found, it runs BuildStatic and returns true so the caller should exit.
// Designed to be called early in main(), after RegisterHandlers.
//
// Example usage in main.go:
//
//	site.RegisterHandlers(myModule)
//	if site.AutoBuild() {
//	    return
//	}
func AutoBuild() bool {
	for i, arg := range os.Args {
		if arg == staticBuildFlag && i+1 < len(os.Args) {
			outputDir := os.Args[i+1]
			if err := BuildStatic(outputDir); err != nil {
				fmt.Println("ssr-static-build error:", err)
				os.Exit(1)
			}
			return true
		}
	}
	return false
}
