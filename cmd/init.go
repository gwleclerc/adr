/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/gwleclerc/adr/constants"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [flags] <directory>",
	Short: "Initialize ADRs configuration",
	Long: fmt.Sprintf(
		`
Initializes the ADR configuration with a base directory.
This is a a prerequisite to running any other subcommand.
The path to the base directory will be stored in a %s file.`,
		ConfigurationFile,
	),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Printf("%s %s %s\n", Red("invalid argument: please specify a"), RedUnderline("directory"), Red("as first argument."))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		path := filepath.Join(".", args[0])
		fmt.Println(Green("Initialize ADRs configuration"))
		if err := initConfiguration(path); err != nil {
			fmt.Println(Red("unable to init ADRs directory: %v", err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		fmt.Println(Green("ADRs configuration has been successfully initialized at %q", path))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initConfiguration(path string) error {
	fmt.Printf("Test if %q directory is already exists\n", path)
	info, err := os.Stat(path)
	switch {
	case err == nil && !info.IsDir():
		return fmt.Errorf("%q is not a directory", path)
	case err != nil && !os.IsNotExist(err):
		return err
	}
	fmt.Printf("Create %q directory because it does not exists\n", path)
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	fmt.Printf("Directory %q has been successfully created\n", path)

	fmt.Printf("Store directory path into %q configuration file\n", ConfigurationFile)
	f, err := os.Create(ConfigurationFile)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := yaml.Marshal(Config{
		Directory: path,
	})
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return err
	}
	fmt.Printf("Directory path has been successfully stored in %q\n", ConfigurationFile)
	return nil
}
