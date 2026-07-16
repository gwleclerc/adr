package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/urfave/cli/v3"
)

var (
	// Exit is overridable so tests can assert on the exit code.
	Exit = os.Exit
	// errSilent is returned by actions that already printed their own message.
	errSilent = errors.New("silent error")
)

// BuildInfo carries the version metadata injected at build time.
type BuildInfo struct {
	AppName string
	Version string
	Commit  string
	Date    string
}

// newApp wires the root command and its subcommands.
func newApp(bi BuildInfo) *cli.Command {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("%s version %s\ncommit: %s\nbuilt at: %s\n", bi.AppName, bi.Version, bi.Commit, bi.Date)
	}

	return &cli.Command{
		Name:                  bi.AppName,
		Usage:                 "A tool to manage Architecture Decision Records (ADRs)",
		Version:               bi.Version,
		HideHelpCommand:       true,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			initCommand(),
			newCommand(),
			addCommand(),
			updateCommand(),
			deprecateCommand(),
			supersedeCommand(),
			listCommand(),
			showCommand(),
			editCommand(),
			tocCommand(),
			lintCommand(),
			templateCommand(),
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				fmt.Fprintf(os.Stderr, "%s %s %s\n",
					cs.Red("invalid argument: please specify a"),
					cs.RedUnderline("command"),
					cs.Red("to execute"),
				)
				return errSilent
			}
			printError("unknown command %q", cmd.Args().First())
			return errSilent
		},
		// Actions print their own (colored) messages, so suppress the default
		// handler to avoid printing errors twice.
		ExitErrHandler: func(_ context.Context, _ *cli.Command, _ error) {},
	}
}

// Execute runs the CLI. It is called by main.main().
func Execute(bi BuildInfo) {
	if err := newApp(bi).Run(context.Background(), os.Args); err != nil {
		Exit(1)
	}
}

// missingArgument prints the canonical "please specify a <what> in arguments" error.
func missingArgument(what string) {
	fmt.Fprintf(os.Stderr, "%s %s %s\n",
		cs.Red("invalid argument: please specify a"),
		cs.RedUnderline(what),
		cs.Red("in arguments"),
	)
}

// printError prints a red, formatted error message to stderr.
func printError(format string, a ...any) {
	fmt.Fprintln(os.Stderr, cs.Red(format, a...))
}

// printWarning prints a yellow, formatted warning message to stderr.
func printWarning(format string, a ...any) {
	fmt.Fprintln(os.Stderr, cs.Yellow(format, a...))
}

// splitCSV flattens comma-separated flag values into a clean list of tokens,
// so both `-t a,b` and `-t a -t b` yield the same result.
func splitCSV(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			if part = strings.TrimSpace(part); part != "" {
				out = append(out, part)
			}
		}
	}
	return out
}
