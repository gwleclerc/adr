package cmd

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/gwleclerc/adr/utils"
	"github.com/urfave/cli/v3"
)

type lintIssue struct {
	File    string `json:"file"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func lintCommand() *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "Check the ADRs for consistency problems",
		Description: `Report inconsistencies across records: dangling superseder references,
duplicate numbers, invalid statuses, superseders without a superseded status, and missing
titles. Exits non-zero when any issue is found (useful in CI).`,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "json", Usage: "output issues as JSON"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}
			issues := lintRecords(service.GetRecords())

			if cmd.Bool("json") {
				if err := printJSON(issues); err != nil {
					printError("unable to encode issues: %v", err)
					return errSilent
				}
			} else if len(issues) == 0 {
				fmt.Println(cs.Green("No issues found."))
			} else {
				for _, is := range issues {
					fmt.Println(cs.Yellow("%s: %s (%s)", is.File, is.Message, is.Rule))
				}
			}
			if len(issues) > 0 {
				return errSilent
			}
			return nil
		},
	}
}

// lintRecords returns every consistency problem found across the records.
func lintRecords(adrs []records.AdrData) []lintIssue {
	ids := make(map[string]bool, len(adrs))
	for _, a := range adrs {
		ids[a.ID] = true
	}

	numbers := map[string][]string{}
	issues := []lintIssue{}
	for _, a := range adrs {
		if a.Title == "" {
			issues = append(issues, lintIssue{a.Name, "missing-title", "record has no title"})
		}
		if !slices.Contains(records.AdrStatuses, a.Status) {
			issues = append(issues, lintIssue{a.Name, "invalid-status", fmt.Sprintf("unknown status %q", a.Status)})
		}
		for superseder := range a.Superseders {
			if !ids[superseder] {
				issues = append(issues, lintIssue{a.Name, "dangling-superseder", fmt.Sprintf("superseder %q does not exist", superseder)})
			}
		}
		if len(a.Superseders) > 0 && a.Status != records.SUPERSEDED {
			issues = append(issues, lintIssue{a.Name, "inconsistent-status", fmt.Sprintf("has superseders but status is %q, not superseded", a.Status)})
		}
		if number := utils.GetRecordNumber(a.Name); number != "" {
			numbers[number] = append(numbers[number], a.Name)
		}
	}
	for number, files := range numbers {
		if len(files) > 1 {
			sort.Strings(files)
			issues = append(issues, lintIssue{files[0], "duplicate-number", fmt.Sprintf("number %s is used by %s", number, strings.Join(files, ", "))})
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].File != issues[j].File {
			return issues[i].File < issues[j].File
		}
		return issues[i].Rule < issues[j].Rule
	})
	return issues
}
