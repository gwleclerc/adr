package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func initCommand() *cli.Command {
	return &cli.Command{
		Name:      "init",
		Usage:     "Initialize ADRs configuration",
		ArgsUsage: "<directory>",
		Description: fmt.Sprintf(`Initializes the ADR configuration with a base directory.
This is a prerequisite to running any other subcommand.
The path to the base directory will be stored in a %s file.`, cs.ConfigurationFile),
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				missingArgument("directory")
				return errSilent
			}
			path := filepath.Join(".", cmd.Args().First())
			if err := initConfiguration(path); err != nil {
				printError("unable to init ADRs directory: %v", err)
				return errSilent
			}
			return nil
		},
	}
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

	b, err := yaml.Marshal(cs.Config{
		Directory: path,
	})
	if err != nil {
		return err
	}
	if err := os.WriteFile(cs.ConfigurationFile, b, 0o644); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(cs.Green("ADRs configuration has been successfully initialized at %q", path))
	fmt.Println()
	return nil
}
