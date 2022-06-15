package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [flags] <directory>",
	Short: "Initialize ADRs configuration",
	Long: fmt.Sprintf(`
Initializes the ADR configuration with a base directory.
This is a a prerequisite to running any other subcommand.
The path to the base directory will be stored in a %s file.`,
		cs.ConfigurationFile,
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			fmt.Printf("%s %s %s\n", cs.Red("invalid argument: please specify a"), cs.RedUnderline("directory"), cs.Red("in arguments"))
			return ErrSilent
		}
		path := filepath.Join(".", args[0])
		if err := initConfiguration(path); err != nil {
			fmt.Println(cs.Red("unable to init ADRs directory: %v", err))
			return ErrSilent
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initConfiguration(path string) error {
	info, err := os.Stat(path)

	switch {
	case err == nil && !info.IsDir():
		return fmt.Errorf("%q is not a directory", path)
	case err != nil && !os.IsNotExist(err):
		return err
	}

	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(cs.ConfigurationFile)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := yaml.Marshal(cs.Config{
		Directory: path,
	})
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(cs.Green("ADRs configuration has been successfully initialized at %q", path))
	fmt.Println()
	return nil
}
