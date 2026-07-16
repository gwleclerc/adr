package records

import (
	"strings"
	"testing"
)

func TestParseStatus(t *testing.T) {
	for _, s := range AdrStatuses {
		got, err := ParseStatus(string(s))
		if err != nil {
			t.Errorf("ParseStatus(%q) returned error: %v", s, err)
		}
		if got != s {
			t.Errorf("ParseStatus(%q) = %q, want %q", s, got, s)
		}
	}
	if _, err := ParseStatus("bogus"); err == nil {
		t.Error("ParseStatus(\"bogus\") should return an error")
	}
}

func TestAllowedStatuses(t *testing.T) {
	got := AllowedStatuses()
	for _, s := range AdrStatuses {
		if !strings.Contains(got, string(s)) {
			t.Errorf("AllowedStatuses() = %q, missing %q", got, s)
		}
	}
	if !strings.Contains(got, " or ") {
		t.Errorf("AllowedStatuses() = %q, expected an \" or \" separator", got)
	}
}
