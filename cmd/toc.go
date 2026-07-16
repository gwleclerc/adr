package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/gwleclerc/adr/utils"
	"github.com/urfave/cli/v3"
)

func tocCommand() *cli.Command {
	return &cli.Command{
		Name:  "toc",
		Usage: "Generate a table of contents for the ADRs",
		Description: `Generate a markdown index of every record (number, title, status, date).
Writes to stdout by default, or to a file with --output (e.g. docs/adrs/README.md).`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "write the index to a file instead of stdout",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}
			toc := renderTOC(service.GetRecords())
			if out := cmd.String("output"); out != "" {
				if err := os.WriteFile(out, []byte(toc), 0o644); err != nil {
					printError("unable to write %q: %v", out, err)
					return errSilent
				}
				fmt.Println(cs.Green("Table of contents written to %q", out))
				return nil
			}
			fmt.Print(toc)
			return nil
		},
	}
}

// renderTOC builds a markdown index of the records. Links are the record filenames,
// so the index is meant to live alongside the records (e.g. in the ADR directory).
func renderTOC(adrs []records.AdrData) string {
	var b strings.Builder
	b.WriteString("# Architecture Decision Records\n\n")
	if len(adrs) == 0 {
		b.WriteString("_No records yet._\n")
		return b.String()
	}
	b.WriteString("| # | Title | Status | Date |\n|---|---|---|---|\n")
	for _, a := range adrs {
		number := utils.GetRecordNumber(a.Name)
		if number == "" {
			number = "-"
		}
		date := "-"
		if !a.CreationDate.IsZero() {
			date = a.CreationDate.Format("2006-01-02")
		}
		title := strings.ReplaceAll(a.Title, "|", "\\|")
		fmt.Fprintf(&b, "| %s | [%s](%s) | %s | %s |\n", number, title, a.Name, a.Status, date)
	}
	return b.String()
}
