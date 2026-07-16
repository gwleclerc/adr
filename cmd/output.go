package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/olekukonko/tablewriter"
)

// printJSON writes v to stdout as indented JSON.
func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

// reportRecord prints a just-updated record, as JSON when jsonOut is set,
// otherwise as a confirmation message followed by a one-row table.
func reportRecord(record records.AdrData, jsonOut bool) error {
	if jsonOut {
		return printJSON(record)
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
