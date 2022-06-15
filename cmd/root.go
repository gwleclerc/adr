package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/spf13/cobra"
)

var (
	Exit      = os.Exit
	ErrSilent = errors.New("SilentErr")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "adr",
	Short:         "A tool to manage Architecture Decision Records (ADRs)",
	Long:          `A tool to manage Architecture Decision Records (ADRs)`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s %s %s\n", cs.Red("invalid argument: please specify a"), cs.RedUnderline("command"), cs.Red("to execute"))
		return ErrSilent
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		Exit(1)
	}
}

func init() {
	cobra.AddTemplateFunc("StyleHeading", cs.Green)
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
	).Replace(usageTemplate)

	// To avoid conflicts with 'Global Flags' we use regex for 'Flags'
	re := regexp.MustCompile(`(?m)^Flags:\s*$`)
	usageTemplate = re.ReplaceAllLiteralString(usageTemplate, `{{StyleHeading "Flags:"}}`)

	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		fmt.Println(cs.Red(err.Error()))
		return ErrSilent
	})
}
