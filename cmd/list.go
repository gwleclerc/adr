package cmd

import (
	"fmt"
	"os"

	. "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	list_authors []string
	list_tags    []string
	list_status  []string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List ADR files",
	Long: fmt.Sprintf(
		`
List ADR files present in directory stored in %s configuration file.`,
		ConfigurationFile,
	),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := utils.RetrieveADRsPath()
		if err != nil {
			fmt.Println(Red("unable to retrieve ADRs path, you should look at the %s configuration file: %v", ConfigurationFile, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if err := listRecords(path); err != nil {
			fmt.Println(Red("unable to list ADRs in directory %q: %v", path, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
	},
}

func init() {
	listCmd.Flags().StringSliceVarP(
		&list_authors,
		"authors",
		"a",
		[]string{},
		"filter records by authors",
	)
	listCmd.Flags().StringSliceVarP(
		&list_status,
		"status",
		"s",
		[]string{},
		"filter records by status",
	)
	listCmd.Flags().StringSliceVarP(
		&new_tags,
		"tags",
		"t",
		[]string{},
		"filter records by tags",
	)
	rootCmd.AddCommand(listCmd)
}

func listRecords(path string) error {
	adrs, err := utils.IndexADRs(path)
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(TableHeader)
	data := [][]string{}
	for _, adr := range adrs {
		if len(list_authors) > 0 {
			if !slices.Contains(list_authors, adr.Author) {
				continue
			}
		}
		if len(list_status) > 0 {
			if !slices.Contains(list_status, adr.Status.String()) {
				continue
			}
		}
		if len(list_tags) > 0 {
			found := false
			for tag := range adr.Tags {
				if slices.Contains(list_tags, tag) {
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
