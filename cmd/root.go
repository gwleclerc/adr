/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	. "github.com/gwleclerc/adr/constants"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "adr",
	Short:         "A tool to manage Architecture Decision Records (ADRs)",
	Long:          `A tool to manage Architecture Decision Records (ADRs)`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

var SilentErr = errors.New("SilentErr")

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.AddTemplateFunc("StyleHeading", Green)
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
		fmt.Println(Red(err.Error()))
		fmt.Println(cmd.UsageString())
		return SilentErr
	})
}
