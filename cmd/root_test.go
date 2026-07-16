package cmd

import (
	"reflect"
	"testing"
)

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"comma separated", []string{"a,b"}, []string{"a", "b"}},
		{"repeated flags", []string{"a", "b"}, []string{"a", "b"}},
		{"mixed with spaces", []string{"a, b", "c"}, []string{"a", "b", "c"}},
		{"empty tokens dropped", []string{"", "a,,b", " "}, []string{"a", "b"}},
		{"nothing", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitCSV(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitCSV(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
