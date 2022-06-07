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
	"github.com/gwleclerc/adr/types"
	"github.com/gwleclerc/adr/utils"
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
			fmt.Printf("%s %s %s\n", Red("invalid argument: please specify a"), RedUnderline("record ID"), Red("as arguments"))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if len(args) > 1 {
			fmt.Println(Yellow("too many argument: keeping only the first record ID"))
		}
		recordID := args[0]
		if len(add_tags) <= 0 && len(add_superseders) <= 0 {
			fmt.Println(Red("invalid arguments: nothing to add to the record %q", recordID))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		path, err := utils.RetrieveADRsPath()
		if err != nil {
			fmt.Println(Red("unable to retrieve ADRs path, you should look at the %s configuration file: %v", ConfigurationFile, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if err := addToRecord(path, recordID); err != nil {
			fmt.Println(Red("unable to update ADR %q: %v", recordID, err))
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

func addToRecord(path, recordID string) error {
	adrs, err := utils.IndexADRs(path)
	if err != nil {
		return err
	}

	var record *types.AdrData
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
	if len(add_tags) > 0 {
		record.Tags.Append(add_tags...)
	}
	if len(add_superseders) > 0 {
		record.Superseders.Append(add_superseders...)
	}
	record.LastUpdateDate = time.Now()

	b, err := types.MarshalYAML(record)
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
