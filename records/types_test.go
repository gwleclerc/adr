package records

import (
	"reflect"
	"testing"
)

func TestSetAppend(t *testing.T) {
	type args struct {
		set              Set[string]
		elementsToAppend []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "nil", args: args{nil, []string{"elem"}}, want: []string{"elem"}},
		{name: "empty", args: args{Set[string]{}, []string{"elem"}}, want: []string{"elem"}},
		{name: "not empty", args: args{Set[string]{"elem1": true}, []string{"elem2"}}, want: []string{"elem1", "elem2"}},
		{name: "several", args: args{Set[string]{"elem1": true}, []string{"elem2", "elem3"}}, want: []string{"elem1", "elem2", "elem3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.set.Append(tt.args.elementsToAppend...)
			got := tt.args.set.ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdrStatusCompletion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetRemove(t *testing.T) {
	type args struct {
		set              Set[string]
		elementsToRemove []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "nil", args: args{nil, []string{"elem"}}, want: []string{}},
		{name: "empty", args: args{Set[string]{}, []string{"elem"}}, want: []string{}},
		{name: "not empty", args: args{Set[string]{"elem1": true}, []string{"elem1"}}, want: []string{}},
		{name: "several", args: args{Set[string]{"elem1": true, "elem2": true}, []string{"elem1"}}, want: []string{"elem2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.set.Remove(tt.args.elementsToRemove...)
			got := tt.args.set.ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdrStatusCompletion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSet(t *testing.T) {
	type args struct {
		set           Set[string]
		elementsToSet []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "nil", args: args{nil, []string{"elem"}}, want: []string{"elem"}},
		{name: "empty", args: args{Set[string]{}, []string{"elem"}}, want: []string{"elem"}},
		{name: "not empty", args: args{Set[string]{"elem1": true}, []string{"elem2"}}, want: []string{"elem2"}},
		{name: "several", args: args{Set[string]{"elem1": true, "elem2": true}, []string{"elem1"}}, want: []string{"elem1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.set.Set(tt.args.elementsToSet...)
			got := tt.args.set.ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdrStatusCompletion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
