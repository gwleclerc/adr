package constants

import (
	"time"

	"github.com/fatih/color"
)

const (
	ConfigurationFile = ".adrrc"
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
	PROPOSED   AdrStatus = "Proposed"
	ACCEPTED   AdrStatus = "Accepted"
	DEPRECATED AdrStatus = "Deprecated"
	SUPERSEDED AdrStatus = "Superseded"
)

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
