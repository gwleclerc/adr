package main

import (
	"flag"
	"testing"

	"github.com/gwleclerc/adr/cmd"
	"github.com/spf13/pflag"
)

func TestMain(t *testing.T) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	cmd.Exit = func(code int) {
		t.Errorf("exited with code: %d", code)
	}
	cmd.Execute()
}
