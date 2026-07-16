package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gwleclerc/adr/records"
	"github.com/gwleclerc/adr/templates"
	"github.com/urfave/cli/v3"
)

func templateCommand() *cli.Command {
	return &cli.Command{
		Name:  "template",
		Usage: "Inspect the ADR body templates",
		Description: `List the available ADR body templates and print their contract
(the sections and their guidance) so a body can be authored to match.

Custom templates are picked up from the "templates_dir" declared in the nearest
.adrrc.yml (each *.tpl file becomes a template named after the file).`,
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List the available templates",
				Action: func(_ context.Context, _ *cli.Command) error {
					reg, err := loadTemplates()
					if err != nil {
						printError("unable to load templates: %v", err)
						return errSilent
					}
					for _, name := range templates.Names(reg) {
						source := "custom"
						if reg[name].Builtin {
							source = "built-in"
						}
						fmt.Printf("%-16s (%s)\n", name, source)
					}
					return nil
				},
			},
			{
				Name:      "show",
				Usage:     "Print a template's sections and guidance",
				ArgsUsage: "<template name>",
				Action: func(_ context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() == 0 {
						missingArgument("template name")
						return errSilent
					}
					name := cmd.Args().First()
					reg, err := loadTemplates()
					if err != nil {
						printError("unable to load templates: %v", err)
						return errSilent
					}
					tpl, ok := reg[name]
					if !ok {
						printError("unknown template %q: available: %s", name, strings.Join(templates.Names(reg), ", "))
						return errSilent
					}
					fmt.Print(tpl.Body)
					return nil
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// `adr template` with no subcommand behaves like `adr template list`.
			return cmd.Command("list").Run(ctx, []string{"list"})
		},
	}
}

// loadTemplates loads the template registry, tolerating a missing config so the
// built-ins are still listed outside an initialized project.
func loadTemplates() (map[string]templates.Template, error) {
	dir := ""
	if cfg, base, err := records.LoadConfig(); err == nil && cfg.TemplatesDir != "" {
		dir = filepath.Join(base, cfg.TemplatesDir)
	}
	return templates.Load(dir)
}
