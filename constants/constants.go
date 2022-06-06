package constants

import (
	"fmt"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/spf13/cobra"
)

const (
	ConfigurationFile = ".adrrc.yml"
	DefaultUserName   = "Unknown"
	CreateADRTemplate = "create_adr.tpl"
)

var (
	Red          = gchalk.WithRed().Sprintf
	RedUnderline = gchalk.WithRed().WithUnderline().Sprintf
	Green        = gchalk.WithGreen().Sprintf
	Yellow       = gchalk.WithYellow().Sprintf
	Grey         = gchalk.WithGrey().Sprintf
	White        = gchalk.WithWhite().Sprintf
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

// Colorized returns the AdrStatus as a colored string
func (e AdrStatus) Colorized() string {
	switch e {
	case PROPOSED:
		return Yellow(e.String())
	case ACCEPTED, SUPERSEDED:
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
	}, cobra.ShellCompDirectiveDefault
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
