//go:build integration
// +build integration

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/gwleclerc/adr/cmd"
)

func TestMain(t *testing.T) {
	// The `go test` framework already consumed its own `-test.*` flags at
	// startup; strip them from os.Args so the CLI parser does not choke on them.
	os.Args = stripTestFlags(os.Args)
	cmd.Exit = func(code int) {
		t.Errorf("exited with code: %d", code)
	}
	cmd.Execute(cmd.BuildInfo{
		AppName: "adr",
		Version: "test",
		Commit:  "test",
		Date:    "test",
	})
}

func stripTestFlags(args []string) []string {
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-test.") || strings.HasPrefix(arg, "--test.") {
			// Skip a space-separated value (e.g. `--test.coverprofile out.cov`).
			if !strings.Contains(arg, "=") && i+1 < len(args) {
				i++
			}
			continue
		}
		out = append(out, arg)
	}
	return out
}
