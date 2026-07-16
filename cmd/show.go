package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gwleclerc/adr/records"
	"github.com/urfave/cli/v3"
)

func showCommand() *cli.Command {
	return &cli.Command{
		Name:      "show",
		Usage:     "Print a single ADR",
		ArgsUsage: "<record ID>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "print the record metadata as JSON instead of the file",
			},
		},
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

			if cmd.Bool("json") {
				b, err := json.MarshalIndent(record, "", "  ")
				if err != nil {
					printError("unable to encode record: %v", err)
					return errSilent
				}
				fmt.Println(string(b))
				return nil
			}

			b, err := os.ReadFile(service.RecordPath(record))
			if err != nil {
				printError("unable to read record: %v", err)
				return errSilent
			}
			fmt.Print(string(b))
			return nil
		},
	}
}
