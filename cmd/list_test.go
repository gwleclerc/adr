package cmd

import (
	"reflect"
	"testing"

	"github.com/gwleclerc/adr/records"
)

func mkRecord(id, author string, status records.AdrStatus, tags ...string) records.AdrData {
	set := make(records.Set[string])
	set.Append(tags...)
	return records.AdrData{ID: id, Author: author, Status: status, Tags: set}
}

func TestFilterRecords(t *testing.T) {
	adrs := []records.AdrData{
		mkRecord("1", "alice", records.ACCEPTED, "api"),
		mkRecord("2", "bob", records.DEPRECATED, "db"),
		mkRecord("3", "alice", records.ACCEPTED, "db", "api"),
	}

	ids := func(rs []records.AdrData) []string {
		out := make([]string, 0, len(rs))
		for _, r := range rs {
			out = append(out, r.ID)
		}
		return out
	}

	tests := []struct {
		name    string
		filters listFilters
		want    []string
	}{
		{"no filter", listFilters{}, []string{"1", "2", "3"}},
		{"by author", listFilters{authors: []string{"alice"}}, []string{"1", "3"}},
		{"by status", listFilters{status: []string{"deprecated"}}, []string{"2"}},
		{"by tag", listFilters{tags: []string{"db"}}, []string{"2", "3"}},
		{"author AND tag", listFilters{authors: []string{"alice"}, tags: []string{"db"}}, []string{"3"}},
		{"no match", listFilters{authors: []string{"carol"}}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ids(filterRecords(adrs, tt.filters))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterRecords(%+v) = %v, want %v", tt.filters, got, tt.want)
			}
		})
	}
}
