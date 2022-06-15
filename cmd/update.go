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
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			fmt.Printf("%s %s %s\n", cs.Red("invalid argument: please specify a"), cs.RedUnderline("record ID"), cs.Red("in arguments"))
			return ErrSilent
		}
		if len(args) > 1 {
			fmt.Println(cs.Yellow("too many argument: keeping only the first record ID"))
		}
		recordID := args[0]
		service, err := records.NewService()
		if err != nil {
			fmt.Println(cs.Red("unable to initialize records service: %v", err))
			return ErrSilent
		}
		if err := updateRecord(service, recordID); err != nil {
			fmt.Println(cs.Red("unable to update ADR %q: %v", recordID, err))
			return ErrSilent
		}
		return nil
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
		nil, // must be nil to allow '--tags=' to remove all tags on record
		`tags of the record`,
	)
	updateCmd.Flags().StringSliceVarP(
		&update_superseders,
		"superseders",
		"r",
		nil, // must be nil to allow '--superseders=' to remove all superseders on record
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
	if update_tags != nil {
		record.Tags.Set(update_tags...)
	}
	if update_superseders != nil {
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
