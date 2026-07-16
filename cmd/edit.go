package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/gwleclerc/adr/records"
	"github.com/urfave/cli/v3"
)

func editCommand() *cli.Command {
	return &cli.Command{
		Name:      "edit",
		Usage:     "Open an ADR in $EDITOR",
		ArgsUsage: "<record ID>",
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				missingArgument("record ID")
				return errSilent
			}
			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}
			record, ok := service.GetRecord(cmd.Args().First())
			if !ok {
				printError("record %q not found", cmd.Args().First())
				return errSilent
			}
			if err := openEditor(service.RecordPath(record)); err != nil {
				printError("unable to open editor: %v", err)
				return errSilent
			}
			return nil
		},
	}
}

// openEditor opens path in the user's editor ($EDITOR, $VISUAL, or vi).
func openEditor(path string) error {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}
	c := exec.Command(editor, path)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s: %w", editor, err)
	}
	return nil
}
