package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
)

func addCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add tags or superseders into ADR",
		ArgsUsage: "<record ID>",
		Description: `Add tags or superseders to an existing architecture decision record.
It will keep the content and only modify the metadata.`,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "tags of the record",
			},
			&cli.StringSliceFlag{
				Name:    "superseders",
				Aliases: []string{"r"},
				Usage:   "superseders of the record",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				missingArgument("record ID")
				return errSilent
			}
			if cmd.Args().Len() > 1 {
				printWarning("too many arguments: keeping only the first record ID")
			}
			recordID := cmd.Args().First()
			tags := splitCSV(cmd.StringSlice("tags"))
			superseders := splitCSV(cmd.StringSlice("superseders"))
			if len(tags) == 0 && len(superseders) == 0 {
				printError("invalid arguments: nothing to add to the record %q", recordID)
				return errSilent
			}
			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}
			if err := addToRecord(service, recordID, tags, superseders); err != nil {
				printError("unable to update ADR %q: %v", recordID, err)
				return errSilent
			}
			return nil
		},
	}
}

func addToRecord(service *records.Service, recordID string, tags, superseders []string) error {
	record, ok := service.GetRecord(recordID)
	if !ok {
		return errors.New("record not found")
	}
	if len(tags) > 0 {
		record.Tags.Append(tags...)
	}
	if len(superseders) > 0 {
		record.Superseders.Append(superseders...)
	}
	if err := service.UpdateRecord(record); err != nil {
		return err
	}
	fmt.Println()
	fmt.Println(cs.Green("Record %q has been successfully updated:", record.ID))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cs.TableHeader)
	table.Append(record.ToRow())
	table.Render()
	fmt.Println()
	return nil
}
