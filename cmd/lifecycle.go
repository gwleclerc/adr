package cmd

import (
	"context"

	"github.com/gwleclerc/adr/records"
	"github.com/urfave/cli/v3"
)

func deprecateCommand() *cli.Command {
	return &cli.Command{
		Name:      "deprecate",
		Usage:     "Mark an ADR as deprecated",
		ArgsUsage: "<record ID>",
		Flags:     []cli.Flag{&cli.BoolFlag{Name: "json", Usage: "print the updated record as JSON"}},
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
			record.Status = records.DEPRECATED
			if err := service.UpdateRecord(record); err != nil {
				printError("unable to update ADR %q: %v", record.ID, err)
				return errSilent
			}
			if err := reportRecord(record, cmd.Bool("json")); err != nil {
				printError("unable to encode record: %v", err)
				return errSilent
			}
			return nil
		},
	}
}

func supersedeCommand() *cli.Command {
	return &cli.Command{
		Name:      "supersede",
		Usage:     "Mark an ADR as superseded by another",
		ArgsUsage: "<superseded ID> <superseder ID>",
		Flags:     []cli.Flag{&cli.BoolFlag{Name: "json", Usage: "print the updated record as JSON"}},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 2 {
				printError("supersede requires <superseded ID> and <superseder ID>")
				return errSilent
			}
			supersededID, supersederID := cmd.Args().Get(0), cmd.Args().Get(1)
			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}
			record, ok := service.GetRecord(supersededID)
			if !ok {
				printError("record %q not found", supersededID)
				return errSilent
			}
			if _, ok := service.GetRecord(supersederID); !ok {
				printWarning("superseder %q does not match any existing record", supersederID)
			}
			record.Status = records.SUPERSEDED
			record.Superseders.Append(supersederID)
			if err := service.UpdateRecord(record); err != nil {
				printError("unable to update ADR %q: %v", record.ID, err)
				return errSilent
			}
			if err := reportRecord(record, cmd.Bool("json")); err != nil {
				printError("unable to encode record: %v", err)
				return errSilent
			}
			return nil
		},
	}
}
