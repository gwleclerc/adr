package main

import "github.com/gwleclerc/adr/cmd"

// Build metadata injected at link time via -ldflags (see Makefile / .goreleaser.yaml).
var (
	appName      = "adr"
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

func main() {
	cmd.Execute(cmd.BuildInfo{
		AppName: appName,
		Version: buildVersion,
		Commit:  buildCommit,
		Date:    buildDate,
	})
}
