package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List ADR files",
		Description: fmt.Sprintf(
			"List ADR files present in directory stored in %s configuration file.",
			cs.ConfigurationFile,
		),
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "authors",
				Aliases: []string{"a"},
				Usage:   "filter records by authors",
			},
			&cli.StringSliceFlag{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "filter records by status",
			},
			&cli.StringSliceFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "filter records by tags",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}
			filters := listFilters{
				authors: splitCSV(cmd.StringSlice("authors")),
				status:  splitCSV(cmd.StringSlice("status")),
				tags:    splitCSV(cmd.StringSlice("tags")),
			}
			if err := listRecords(service, filters); err != nil {
				printError("unable to list ADRs: %v", err)
				return errSilent
			}
			return nil
		},
	}
}

type listFilters struct {
	authors []string
	status  []string
	tags    []string
}

func listRecords(service *records.Service, filters listFilters) error {
	adrs := service.GetRecords()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cs.TableHeader)
	data := [][]string{}

	for _, adr := range adrs {
		if len(filters.authors) > 0 && !slices.Contains(filters.authors, adr.Author) {
			continue
		}
		if len(filters.status) > 0 && !slices.Contains(filters.status, adr.Status.String()) {
			continue
		}
		if len(filters.tags) > 0 {
			found := false
			for tag := range adr.Tags {
				if slices.Contains(filters.tags, tag) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		data = append(data, adr.ToRow())
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
	return nil
}
