package constants

import (
	"github.com/jwalton/gchalk"
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
