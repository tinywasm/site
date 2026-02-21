// sitebuild compiles a tinywasm/site project and generates static HTML/CSS/JS/SVG output.
//
// Usage:
//
//	sitebuild [--out <dir>] <package-path>
//
// The package at <package-path> must call site.AutoBuild() early in its main().
// sitebuild compiles the package, runs it with --ssr-static-build <dir>,
// and exits with the same code as the subprocess.
//
// Example:
//
//	sitebuild --out dist/ ./cmd/myapp
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// builder is the interface for compiling and running a Go binary.
// Defined for testability — the real implementation uses os/exec.
type builder interface {
	Build(pkg, outBin string) error
	Run(bin string, args []string) error
}

// realBuilder is the production implementation of builder.
type realBuilder struct{}

func (r *realBuilder) Build(pkg, outBin string) error {
	cmd := exec.Command("go", "build", "-o", outBin, pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *realBuilder) Run(bin string, args []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func run(b builder, args []string) int {
	outDir := "dist"
	pkg := ""

	// Parse args: sitebuild [--out <dir>] <package-path>
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--out", "-out", "-o":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "sitebuild: --out requires a directory argument")
				return 1
			}
			i++
			outDir = args[i]
		default:
			if pkg != "" {
				fmt.Fprintln(os.Stderr, "sitebuild: unexpected argument:", args[i])
				return 1
			}
			pkg = args[i]
		}
	}

	if pkg == "" {
		fmt.Fprintln(os.Stderr, "Usage: sitebuild [--out <dir>] <package-path>")
		fmt.Fprintln(os.Stderr, "Example: sitebuild --out dist/ ./cmd/myapp")
		return 1
	}

	// Create a temp binary path
	tmpDir, err := os.MkdirTemp("", "sitebuild-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "sitebuild: failed to create temp dir:", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	outBin := filepath.Join(tmpDir, "sitebuild_app")

	// Step 1: compile
	fmt.Println("sitebuild: compiling", pkg)
	if err := b.Build(pkg, outBin); err != nil {
		fmt.Fprintln(os.Stderr, "sitebuild: build failed:", err)
		return 1
	}

	// Step 2: run with --ssr-static-build flag
	fmt.Println("sitebuild: generating static site to", outDir)
	if err := b.Run(outBin, []string{"--ssr-static-build", outDir}); err != nil {
		fmt.Fprintln(os.Stderr, "sitebuild: static build failed:", err)
		return 1
	}

	fmt.Println("sitebuild: done →", outDir)
	return 0
}

func main() {
	os.Exit(run(&realBuilder{}, os.Args[1:]))
}
