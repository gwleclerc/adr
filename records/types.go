package records

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	cs "github.com/gwleclerc/adr/constants"
	"gopkg.in/yaml.v3"
)

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

// AdrStatuses lists every allowed status in display order.
var AdrStatuses = []AdrStatus{UNKNOWN, PROPOSED, ACCEPTED, DEPRECATED, SUPERSEDED, OBSERVED}

// AdrStatusDescriptions documents the meaning of each status.
var AdrStatusDescriptions = map[AdrStatus]string{
	UNKNOWN:    "status is not determined",
	PROPOSED:   "the record has been proposed but is not accepted yet by stakeholders",
	ACCEPTED:   "the record has been accepted by stakeholders",
	DEPRECATED: "the decision record is deprecated and no longer applies",
	SUPERSEDED: "the decision record has been superseded by a new one",
	OBSERVED:   "documents a pre-existing decision reconstructed after the fact, e.g. while making sense of legacy code you did not write",
}

// StatusHelp returns a multi-line block listing each status and its meaning,
// suitable for appending to a command's help description.
func StatusHelp() string {
	var b strings.Builder
	b.WriteString("Statuses:")
	for _, s := range AdrStatuses {
		fmt.Fprintf(&b, "\n  %-11s %s", s, AdrStatusDescriptions[s])
	}
	return b.String()
}

// String is used by fmt.Print and everywhere a status is rendered as text.
func (e AdrStatus) String() string {
	return string(e)
}

// ParseStatus validates a raw value and returns the matching AdrStatus.
func ParseStatus(v string) (AdrStatus, error) {
	status := AdrStatus(v)
	if slices.Contains(AdrStatuses, status) {
		return status, nil
	}
	return "", fmt.Errorf("must be one of %s", AllowedStatuses())
}

// AllowedStatuses returns the allowed statuses formatted for help/error messages,
// e.g. `"unknown", "proposed", "accepted", "deprecated", "superseded" or "observed"`.
func AllowedStatuses() string {
	quoted := make([]string, len(AdrStatuses))
	for i, s := range AdrStatuses {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	if len(quoted) < 2 {
		return strings.Join(quoted, "")
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + " or " + quoted[len(quoted)-1]
}

// Colorized returns the AdrStatus as a colored string
func (e AdrStatus) Colorized() string {
	switch e {
	case PROPOSED:
		return cs.Yellow(e.String())
	case ACCEPTED, OBSERVED:
		return cs.Green(e.String())
	default:
		return cs.Grey(e.String())
	}
}

type AdrData struct {
	ID             string      `yaml:"id" json:"id"`
	Title          string      `yaml:"title" json:"title"`
	Author         string      `yaml:"author" json:"author"`
	Status         AdrStatus   `yaml:"status" json:"status"`
	CreationDate   time.Time   `yaml:"creation_date" mapstructure:"creation_date" json:"creation_date"`
	LastUpdateDate time.Time   `yaml:"last_update_date" mapstructure:"last_update_date" json:"last_update_date"`
	Tags           Set[string] `yaml:"tags,omitempty" json:"tags,omitempty"`
	Superseders    Set[string] `yaml:"superseders,omitempty" json:"superseders,omitempty"`

	Name string `yaml:"-" json:"file"`
	Body string `yaml:"-" json:"-"`
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

type Set[T cmp.Ordered] map[T]bool

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
	slices.Sort(res)
	return res
}

func (s Set[T]) MarshalYAML() (interface{}, error) {
	return s.ToSlice(), nil
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
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
