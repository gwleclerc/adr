package constants

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

const (
	ConfigurationFile = ".adrrc.yml"
	DefaultUserName   = "Unknown"
	CreateADRTemplate = "create_adr.tpl"
)

var (
	Red          = color.New(color.FgRed).SprintfFunc()
	RedUnderline = color.New(color.FgRed, color.Underline).SprintfFunc()
	Green        = color.New(color.FgGreen).SprintfFunc()
	Yellow       = color.New(color.FgYellow).SprintfFunc()
)

type Config struct {
	Directory string `yaml:"directory"`
}

// AdrStatus type
type AdrStatus string

// ADR status enums
const (
	PROPOSED   AdrStatus = "proposed"
	ACCEPTED   AdrStatus = "accepted"
	DEPRECATED AdrStatus = "deprecated"
	SUPERSEDED AdrStatus = "superseded"
	UNKNOWN    AdrStatus = "unknown"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *AdrStatus) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *AdrStatus) Set(v string) error {
	switch AdrStatus(v) {
	case UNKNOWN, PROPOSED, ACCEPTED, DEPRECATED, SUPERSEDED:
		*e = AdrStatus(v)
		return nil
	default:
		return fmt.Errorf(
			"must be one of %q, %q, %q, %q or %q",
			UNKNOWN, PROPOSED, ACCEPTED, DEPRECATED, SUPERSEDED,
		)
	}
}

// Type is only used in help text
func (e *AdrStatus) Type() string {
	return "status"
}

type AdrData struct {
	ID     string    `yaml:"id"`
	Title  string    `yaml:"title"`
	Status AdrStatus `yaml:"status"`
	Date   time.Time `yaml:"date"`
	Author string    `yaml:"author"`
	Tags   []string  `yaml:"tags"`

	Name string `yaml:"-"`
	Body string `yaml:"-"`
}
