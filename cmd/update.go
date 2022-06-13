package cmd

import (
	"errors"
	"fmt"
	"os"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	update_author      string
	update_status      records.AdrStatus
	update_tags        []string
	update_superseders []string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [flags] <record ID>",
	Short: "Update an ADR",
	Long: `
Update an existing architecture decision record.
It will keep the content and only modify the metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Printf("%s %s %s\n", cs.Red("invalid argument: please specify a"), cs.RedUnderline("record ID"), cs.Red("as arguments"))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if len(args) > 1 {
			fmt.Println(cs.Yellow("too many argument: keeping only the first record ID"))
		}
		recordID := args[0]
		service, err := records.NewService()
		if err != nil {
			fmt.Println(cs.Red("unable to initialize records service: %v", err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if err := updateRecord(service, recordID); err != nil {
			fmt.Println(cs.Red("unable to update ADR %q: %v", recordID, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
	},
}

func init() {
	updateCmd.Flags().StringVarP(
		&update_author,
		"author",
		"a",
		"",
		"author of the record",
	)
	updateCmd.Flags().VarP(
		&update_status,
		"status",
		"s",
		`status of the record, allowed: "unknown", "proposed", "accepted", "deprecated", "superseded" or "observed"`,
	)
	_ = updateCmd.RegisterFlagCompletionFunc("status", records.AdrStatusCompletion)
	updateCmd.Flags().StringSliceVarP(
		&update_tags,
		"tags",
		"t",
		[]string{},
		`tags of the record`,
	)
	updateCmd.Flags().StringSliceVarP(
		&update_superseders,
		"superseders",
		"r",
		[]string{},
		`superseders of the record`,
	)
	rootCmd.AddCommand(updateCmd)
}

func updateRecord(service *records.Service, recordID string) error {
	record, ok := service.GetRecord(recordID)
	if !ok {
		return errors.New("record not found")
	}

	if update_author != "" {
		record.Author = update_author
	}
	if update_status != "" {
		record.Status = update_status
	}
	if len(update_tags) > 0 {
		record.Tags.Set(update_tags...)
	}
	if len(update_superseders) > 0 {
		record.Superseders.Set(update_superseders...)
	}

	err := service.UpdateRecord(record)
	if err != nil {
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
