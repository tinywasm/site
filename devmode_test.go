//go:build !wasm

package site

import (
	"flag"
	"os"
	"os/exec"
	"testing"
)

func init() {
	// Define the flag to avoid "flag provided but not defined" error in the subprocess test runner
	if flag.Lookup("dev") == nil {
		flag.Bool("dev", false, "enable dev mode")
	}
}

// TestDevModeArgument verifies that passing -dev argument sets DevMode to true.
// It uses a subprocess to simulate a fresh run where init() parses os.Args.
func TestDevModeArgument(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		if !handler.DevMode {
			os.Exit(1)
		}
		os.Exit(0)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDevModeArgument", "-dev")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		t.Fatalf("process ran with -dev but handler.DevMode was false")
	} else if err != nil {
		t.Fatalf("process failed to run: %v", err)
	}
}
