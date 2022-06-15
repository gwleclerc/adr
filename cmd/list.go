package cmd

import (
	"fmt"
	"os"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
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
		cs.ConfigurationFile,
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		service, err := records.NewService()
		if err != nil {
			fmt.Println(cs.Red("unable to initialize records service: %v", err))
			return ErrSilent
		}
		if err := listRecords(service); err != nil {
			fmt.Println(cs.Red("unable to list ADRs: %v", err))
			return ErrSilent
		}
		return nil
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
		&list_tags,
		"tags",
		"t",
		[]string{},
		"filter records by tags",
	)
	rootCmd.AddCommand(listCmd)
}

func listRecords(service *records.Service) error {
	adrs := service.GetRecords()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cs.TableHeader)
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
