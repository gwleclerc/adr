package constants

import (
	"github.com/jwalton/gchalk"
)

type Config struct {
	Directory    string `yaml:"directory"`
	TemplatesDir string `yaml:"templates_dir,omitempty"`
}

const (
	ConfigurationFile = ".adrrc.yml"
	DefaultUserName   = "Unknown"
)

var (
	Red          = gchalk.WithRed().Sprintf
	RedUnderline = gchalk.WithRed().WithUnderline().Sprintf
	Green        = gchalk.WithGreen().Sprintf
	Yellow       = gchalk.WithYellow().Sprintf
	Grey         = gchalk.WithGrey().Sprintf

	TableHeader = []string{"ID", "Title", "Status", "Author", "Creation Date", "Last Update Date", "Superseders", "Tags"}
)
