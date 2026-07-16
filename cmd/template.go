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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "json", Usage: "output templates as JSON"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					reg, err := loadTemplates()
					if err != nil {
						printError("unable to load templates: %v", err)
						return errSilent
					}
					if cmd.Bool("json") {
						infos := make([]templateInfo, 0, len(reg))
						for _, name := range templates.Names(reg) {
							infos = append(infos, templateInfo{Name: name, Builtin: reg[name].Builtin})
						}
						if err := printJSON(infos); err != nil {
							printError("unable to encode templates: %v", err)
							return errSilent
						}
						return nil
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
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "json", Usage: "output the template (with its headings) as JSON"},
				},
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
					if cmd.Bool("json") {
						detail := templateDetail{
							Name:     name,
							Builtin:  tpl.Builtin,
							Headings: templates.Headings(tpl.Body),
							Body:     tpl.Body,
						}
						if err := printJSON(detail); err != nil {
							printError("unable to encode template: %v", err)
							return errSilent
						}
						return nil
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

type templateInfo struct {
	Name    string `json:"name"`
	Builtin bool   `json:"builtin"`
}

type templateDetail struct {
	Name     string   `json:"name"`
	Builtin  bool     `json:"builtin"`
	Headings []string `json:"headings"`
	Body     string   `json:"body"`
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
