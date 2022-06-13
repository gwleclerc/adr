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
	add_tags        []string
	add_superseders []string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [flags] <record ID>",
	Short: "Add tags or superseders into ADR",
	Long: `
Add tags or superseders to an existing architecture decision record.
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
		if len(add_tags) <= 0 && len(add_superseders) <= 0 {
			fmt.Println(cs.Red("invalid arguments: nothing to add to the record %q", recordID))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		service, err := records.NewService()
		if err != nil {
			fmt.Println(cs.Red("unable to initialize records service: %v", err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if err := addToRecord(service, recordID); err != nil {
			fmt.Println(cs.Red("unable to update ADR %q: %v", recordID, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
	},
}

func init() {
	addCmd.Flags().StringSliceVarP(
		&add_tags,
		"tags",
		"t",
		[]string{},
		`tags of the record`,
	)
	addCmd.Flags().StringSliceVarP(
		&add_superseders,
		"superseders",
		"r",
		[]string{},
		`superseders of the record`,
	)
	rootCmd.AddCommand(addCmd)
}

func addToRecord(service *records.Service, recordID string) error {
	record, ok := service.GetRecord(recordID)
	if !ok {
		return errors.New("record not found")
	}
	if len(add_tags) > 0 {
		record.Tags.Append(add_tags...)
	}
	if len(add_superseders) > 0 {
		record.Superseders.Append(add_superseders...)
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
