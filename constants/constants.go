package constants

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jwalton/gchalk"
	"github.com/spf13/cobra"
)

const (
	ConfigurationFile = ".adrrc.yml"
	DefaultUserName   = "Unknown"
	CreateADRTemplate = "create_adr.tpl"
	UpdateADRTemplate = "update_adr.tpl"
)

var (
	Red          = gchalk.WithRed().Sprintf
	RedUnderline = gchalk.WithRed().WithUnderline().Sprintf
	Green        = gchalk.WithGreen().Sprintf
	Yellow       = gchalk.WithYellow().Sprintf
	Grey         = gchalk.WithGrey().Sprintf
	White        = gchalk.WithWhite().Sprintf

	TableHeader = []string{"ID", "Title", "Status", "Author", "Creation Date", "Last Update Date", "Superseders", "Tags"}
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
	ID             string    `yaml:"id"`
	Title          string    `yaml:"title"`
	Author         string    `yaml:"author"`
	Status         AdrStatus `yaml:"status"`
	CreationDate   time.Time `yaml:"creation_date"`
	LastUpdateDate time.Time `yaml:"last_update_date"`
	Tags           []string  `yaml:"tags,omitempty"`
	Superseders    []string  `yaml:"superseders,omitempty"`

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
		strings.Join(a.Superseders, ", "),
		strings.Join(a.Tags, ", "),
	}
}
