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

type updateRecordOptions struct {
	author         string
	status         records.AdrStatus
	setTags        bool
	tags           []string
	setSuperseders bool
	superseders    []string
}

func updateCommand() *cli.Command {
	return &cli.Command{
		Name:      "update",
		Usage:     "Update an ADR",
		ArgsUsage: "<record ID>",
		Description: `Update an existing architecture decision record.
It will keep the content and only modify the metadata.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "author",
				Aliases: []string{"a"},
				Usage:   "author of the record",
			},
			&cli.StringFlag{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "status of the record, allowed: " + records.AllowedStatuses(),
			},
			&cli.StringSliceFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "tags of the record (use --tags= to remove all tags)",
			},
			&cli.StringSliceFlag{
				Name:    "superseders",
				Aliases: []string{"r"},
				Usage:   "superseders of the record (use --superseders= to remove all superseders)",
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

			var status records.AdrStatus
			if cmd.IsSet("status") {
				s, err := records.ParseStatus(cmd.String("status"))
				if err != nil {
					printError("invalid status: %v", err)
					return errSilent
				}
				status = s
			}

			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}
			opts := updateRecordOptions{
				author:         cmd.String("author"),
				status:         status,
				setTags:        cmd.IsSet("tags"),
				tags:           splitCSV(cmd.StringSlice("tags")),
				setSuperseders: cmd.IsSet("superseders"),
				superseders:    splitCSV(cmd.StringSlice("superseders")),
			}
			if err := updateRecord(service, recordID, opts); err != nil {
				printError("unable to update ADR %q: %v", recordID, err)
				return errSilent
			}
			return nil
		},
	}
}

func updateRecord(service *records.Service, recordID string, opts updateRecordOptions) error {
	record, ok := service.GetRecord(recordID)
	if !ok {
		return errors.New("record not found")
	}

	if opts.author != "" {
		record.Author = opts.author
	}
	if opts.status != "" {
		record.Status = opts.status
	}
	if opts.setTags {
		record.Tags.Set(opts.tags...)
	}
	if opts.setSuperseders {
		record.Superseders.Set(opts.superseders...)
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
