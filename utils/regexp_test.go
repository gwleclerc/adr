package utils

import "testing"

func TestGetRecordNumber(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"standard", "001_title.md", "001"},
		{"multi digit", "042_some_title.md", "042"},
		{"no separator", "123.md", ""},
		{"no number", "title.md", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRecordNumber(tt.in); got != tt.want {
				t.Errorf("GetRecordNumber(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
