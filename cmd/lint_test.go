package cmd

import (
	"testing"

	"github.com/gwleclerc/adr/records"
)

func mkFull(name, id, title string, status records.AdrStatus, superseders ...string) records.AdrData {
	set := make(records.Set[string])
	set.Append(superseders...)
	return records.AdrData{Name: name, ID: id, Title: title, Status: status, Superseders: set}
}

func hasRule(issues []lintIssue, rule string) bool {
	for _, i := range issues {
		if i.Rule == rule {
			return true
		}
	}
	return false
}

func TestLintRecordsClean(t *testing.T) {
	clean := []records.AdrData{
		mkFull("001_a.md", "a", "A", records.ACCEPTED),
		mkFull("002_b.md", "b", "B", records.SUPERSEDED, "a"), // superseded by an existing record
	}
	if issues := lintRecords(clean); len(issues) != 0 {
		t.Errorf("expected no issues, got %+v", issues)
	}
}

func TestLintRecordsProblems(t *testing.T) {
	tests := []struct {
		name string
		adrs []records.AdrData
		rule string
	}{
		{
			"dangling superseder",
			[]records.AdrData{mkFull("001_a.md", "a", "A", records.SUPERSEDED, "ghost")},
			"dangling-superseder",
		},
		{
			"superseders without superseded status",
			[]records.AdrData{
				mkFull("001_a.md", "a", "A", records.ACCEPTED, "b"),
				mkFull("002_b.md", "b", "B", records.ACCEPTED),
			},
			"inconsistent-status",
		},
		{
			"duplicate number",
			[]records.AdrData{
				mkFull("001_a.md", "a", "A", records.ACCEPTED),
				mkFull("001_b.md", "b", "B", records.ACCEPTED),
			},
			"duplicate-number",
		},
		{
			"missing title",
			[]records.AdrData{mkFull("001_a.md", "a", "", records.ACCEPTED)},
			"missing-title",
		},
		{
			"invalid status",
			[]records.AdrData{mkFull("001_a.md", "a", "A", records.AdrStatus("bogus"))},
			"invalid-status",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := lintRecords(tt.adrs)
			if !hasRule(issues, tt.rule) {
				t.Errorf("expected rule %q, got %+v", tt.rule, issues)
			}
		})
	}
}
