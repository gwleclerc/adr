package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/gwleclerc/adr/templates"
	"github.com/tcnksm/go-gitconfig"
	"github.com/teris-io/shortid"
	"github.com/urfave/cli/v3"
)

type newRecordOptions struct {
	author     string
	status     records.AdrStatus
	tags       []string
	supersedes []string
	body       string
	edit       bool
}

func newCommand() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a new ADR",
		ArgsUsage: "<record title...>",
		Description: fmt.Sprintf(`Create a new architecture decision record.
It will be created in the directory defined in the nearest %s configuration file.

%s`, cs.ConfigurationFile, records.StatusHelp()),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "author",
				Aliases: []string{"a"},
				Usage:   "author of the record",
			},
			&cli.StringFlag{
				Name:    "status",
				Aliases: []string{"s"},
				Value:   string(records.ACCEPTED),
				Usage:   "status of the record, allowed: " + records.AllowedStatuses(),
			},
			&cli.StringSliceFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "tags of the record",
			},
			&cli.StringSliceFlag{
				Name:    "supersedes",
				Aliases: []string{"r"},
				Usage:   "record ids superseded by this one",
			},
			&cli.StringFlag{
				Name:  "template",
				Value: "bare",
				Usage: "body template name (see `adr template list`)",
			},
			&cli.StringFlag{
				Name:  "body-file",
				Usage: "read the record body from a file (or - for stdin) instead of the template; validated against --template",
			},
			&cli.BoolFlag{
				Name:  "edit",
				Usage: "open the created record in $EDITOR",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			title := strings.TrimSpace(strings.Join(cmd.Args().Slice(), " "))
			if title == "" {
				missingArgument("title")
				return errSilent
			}
			status, err := records.ParseStatus(cmd.String("status"))
			if err != nil {
				printError("invalid status: %v", err)
				return errSilent
			}
			service, err := records.NewService()
			if err != nil {
				printError("unable to initialize records service: %v", err)
				return errSilent
			}

			reg, err := templates.Load(service.TemplatesDir())
			if err != nil {
				printError("unable to load templates: %v", err)
				return errSilent
			}
			templateName := cmd.String("template")
			if !cmd.IsSet("template") && service.DefaultTemplate() != "" {
				templateName = service.DefaultTemplate()
			}
			tpl, ok := reg[templateName]
			if !ok {
				printError("invalid template %q: available: %s", templateName, strings.Join(templates.Names(reg), ", "))
				return errSilent
			}
			body := tpl.Body
			if cmd.IsSet("body-file") {
				content, err := readBody(cmd.String("body-file"))
				if err != nil {
					printError("unable to read body: %v", err)
					return errSilent
				}
				if err := templates.Validate(tpl.Body, content); err != nil {
					printError("invalid body for template %q: %v", templateName, err)
					return errSilent
				}
				body = content
			}

			author := cmd.String("author")
			if author == "" {
				author = service.DefaultAuthor()
			}
			opts := newRecordOptions{
				author:     author,
				status:     status,
				tags:       splitCSV(cmd.StringSlice("tags")),
				supersedes: splitCSV(cmd.StringSlice("supersedes")),
				body:       body,
				edit:       cmd.Bool("edit"),
			}
			if err := newRecord(service, title, opts); err != nil {
				printError("unable to create a new ADR: %v", err)
				return errSilent
			}
			return nil
		},
	}
}

// readBody reads a record body from a file, or from stdin when path is "-".
func readBody(path string) (string, error) {
	if path == "-" {
		b, err := io.ReadAll(os.Stdin)
		return string(b), err
	}
	b, err := os.ReadFile(path)
	return string(b), err
}

func newRecord(service *records.Service, title string, opts newRecordOptions) error {
	author := opts.author
	if author == "" {
		author = resolveAuthor()
	}

	// Since IDs starting with '-' would be interpreted as CLI flags, we regenerate
	// a new ID until this is no longer the case.
	id := shortid.MustGenerate()
	for strings.HasPrefix(id, "-") {
		id = shortid.MustGenerate()
	}

	record := records.AdrData{
		ID:     id,
		Status: opts.status,
		Author: author,
		Tags:   make(records.Set[string]),
	}
	record.Tags.Append(opts.tags...)

	path, err := service.CreateRecord(title, record, opts.body)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(cs.Green("Record has been successfully created with ID %q", record.ID))

	markSuperseded(service, record.ID, opts.supersedes)

	fmt.Println()

	if opts.edit {
		return openEditor(path)
	}
	return nil
}

// resolveAuthor determines the record author from the git config, falling back
// to the OS user and finally to the default user name.
func resolveAuthor() string {
	username, err := gitconfig.Username()
	if err == nil {
		return username
	}
	printWarning("Unable to find a git user: %v", err)
	u, err := user.Current()
	if err != nil {
		printWarning("Unable to find an OS user: %v", err)
		return cs.DefaultUserName
	}
	return u.Username
}
