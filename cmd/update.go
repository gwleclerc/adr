package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/templates"
	"github.com/gwleclerc/adr/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	update_author      string
	update_status      AdrStatus
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
			fmt.Printf("%s %s %s\n", Red("invalid argument: please specify a"), RedUnderline("record ID"), Red("as arguments"))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if len(args) > 1 {
			fmt.Printf(Yellow("too many argument: keeping only the first record ID"))
		}
		recordID := args[0]
		path, err := utils.RetrieveADRsPath()
		if err != nil {
			fmt.Println(Red("unable to retrieve ADRs path, you should look at the %s configuration file: %v", ConfigurationFile, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if err := updateRecord(path, recordID); err != nil {
			fmt.Println(Red("unable to update ADR %q: %v", recordID, err))
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
	updateCmd.RegisterFlagCompletionFunc("status", AdrStatusCompletion)
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

func updateRecord(path, recordID string) error {
	adrs, err := utils.IndexADRs(path)
	if err != nil {
		return err
	}

	var record *AdrData
	for i := range adrs {
		adr := adrs[i]
		if adr.ID == recordID {
			record = &adr
			break
		}
	}
	if record == nil {
		return errors.New("record not found")
	}

	if update_author != "" {
		record.Author = update_author
	}
	if update_status != "" {
		record.Status = update_status
	}
	if len(update_tags) > 0 {
		record.Tags = update_tags
	}
	if len(update_superseders) > 0 {
		record.Superseders = update_superseders
	}
	record.LastUpdateDate = time.Now()

	b, err := utils.MarshalYAML(record)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, record.Name))
	if err != nil {
		return err
	}

	err = templates.Templates[UpdateADRTemplate].Execute(file, map[string]any{
		"Header": strings.Trim(string(b), "\n"),
		"Body":   record.Body,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(Green("Record %q has been successfully updated:", record.ID))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(TableHeader)
	table.Append(record.ToRow())
	table.Render()
	fmt.Println()
	return nil
}
