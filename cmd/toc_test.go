package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/gwleclerc/adr/records"
)

func TestRenderTOCEmpty(t *testing.T) {
	if out := renderTOC(nil); !strings.Contains(out, "No records yet") {
		t.Errorf("empty TOC should note there are no records, got:\n%s", out)
	}
}

func TestRenderTOC(t *testing.T) {
	adrs := []records.AdrData{
		{Name: "001_first.md", Title: "First", Status: records.ACCEPTED, CreationDate: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
		{Name: "002_second.md", Title: "Sec | ond", Status: records.DEPRECATED}, // zero date, pipe in title
	}
	out := renderTOC(adrs)
	for _, want := range []string{
		"# Architecture Decision Records",
		"| 001 | [First](001_first.md) | accepted | 2026-01-02 |",
		"| 002 | [Sec \\| ond](002_second.md) | deprecated | - |", // escaped pipe, unset date
	} {
		if !strings.Contains(out, want) {
			t.Errorf("TOC missing %q\n---\n%s", want, out)
		}
	}
}
