package types

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	. "github.com/gwleclerc/adr/constants"
	"github.com/spf13/cobra"
	"golang.org/x/exp/constraints"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Directory string `yaml:"directory"`
}

// AdrStatus type
type AdrStatus string

// ADR status enums
const (
	UNKNOWN    AdrStatus = "unknown"
	PROPOSED   AdrStatus = "proposed"
	ACCEPTED   AdrStatus = "accepted"
	DEPRECATED AdrStatus = "deprecated"
	SUPERSEDED AdrStatus = "superseded"
	OBSERVED   AdrStatus = "observed"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *AdrStatus) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *AdrStatus) Set(v string) error {
	switch AdrStatus(v) {
	case UNKNOWN, PROPOSED, ACCEPTED, DEPRECATED, SUPERSEDED, OBSERVED:
		*e = AdrStatus(v)
		return nil
	default:
		return fmt.Errorf(
			"must be one of %q, %q, %q, %q, %q or %q",
			UNKNOWN, PROPOSED, ACCEPTED, DEPRECATED, SUPERSEDED, OBSERVED,
		)
	}
}

// Colorized returns the AdrStatus as a colored string
func (e AdrStatus) Colorized() string {
	switch e {
	case PROPOSED:
		return Yellow(e.String())
	case ACCEPTED, OBSERVED:
		return Green(e.String())
	default:
		return Grey(e.String())
	}
}

// Type is only used in help text
func (e *AdrStatus) Type() string {
	return "status"
}

// AdrStatusCompletion is used to autocomplete cobra command
func AdrStatusCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		fmt.Sprintf("%s\t%s", UNKNOWN, "status is not determined"),
		fmt.Sprintf("%s\t%s", PROPOSED, "the record has been proposed but is not accepted yet by stakeholders"),
		fmt.Sprintf("%s\t%s", ACCEPTED, "the record has been accepted by stakeholders"),
		fmt.Sprintf("%s\t%s", DEPRECATED, "the decision record is deprecated and no longer applies"),
		fmt.Sprintf("%s\t%s", SUPERSEDED, "the decision record has been superseded by a new one"),
		fmt.Sprintf("%s\t%s", OBSERVED, "the decision was observed after the fact"),
	}, cobra.ShellCompDirectiveDefault
}

type AdrData struct {
	ID             string      `yaml:"id"`
	Title          string      `yaml:"title"`
	Author         string      `yaml:"author"`
	Status         AdrStatus   `yaml:"status"`
	CreationDate   time.Time   `yaml:"creation_date"`
	LastUpdateDate time.Time   `yaml:"last_update_date"`
	Tags           Set[string] `yaml:"tags,omitempty"`
	Superseders    Set[string] `yaml:"superseders,omitempty"`

	Name string `yaml:"-"`
	Body string `yaml:"-"`
}

func (a AdrData) ToRow() []string {
	return []string{
		a.ID,
		a.Title,
		a.Status.Colorized(),
		a.Author,
		humanize.Time(a.CreationDate),
		humanize.Time(a.LastUpdateDate),
		strings.Join(a.Superseders.ToSlice(), ", "),
		strings.Join(a.Tags.ToSlice(), ", "),
	}
}

type Set[T constraints.Ordered] map[T]bool

func (s *Set[T]) Append(elements ...T) {
	if s == nil || (*s) == nil {
		(*s) = make(Set[T])
	}
	for _, elem := range elements {
		(*s)[elem] = true
	}
}

func (s Set[T]) Remove(elements ...T) {
	for _, elem := range elements {
		delete((s), elem)
	}
}

func (s *Set[T]) Set(elements ...T) {
	(*s) = make(Set[T])
	s.Append(elements...)
}

func (s Set[T]) ToSlice() []T {
	res := make([]T, 0, len(s))
	for elem := range s {
		res = append(res, elem)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res
}

func (s Set[T]) MarshalYAML() (interface{}, error) {
	return s.ToSlice(), nil
}

func (s *Set[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp []T
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	(*s).Set(tmp...)
	return nil
}

func MarshalYAML(v any) ([]byte, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2) // this is what you're looking for
	err := yamlEncoder.Encode(v)
	return b.Bytes(), err
}
