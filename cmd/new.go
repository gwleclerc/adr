package cmd

import (
	"context"
	"fmt"
	"os/user"
	"strings"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/tcnksm/go-gitconfig"
	"github.com/teris-io/shortid"
	"github.com/urfave/cli/v3"
)

type newRecordOptions struct {
	author     string
	status     records.AdrStatus
	tags       []string
	supersedes []string
}

func newCommand() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a new ADR",
		ArgsUsage: "<record title...>",
		Description: fmt.Sprintf(`Create a new architecture decision record.
It will be created in the directory defined in the nearest %s configuration file.`, cs.ConfigurationFile),
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
			opts := newRecordOptions{
				author:     cmd.String("author"),
				status:     status,
				tags:       splitCSV(cmd.StringSlice("tags")),
				supersedes: splitCSV(cmd.StringSlice("supersedes")),
			}
			if err := newRecord(service, title, opts); err != nil {
				printError("unable to create a new ADR: %v", err)
				return errSilent
			}
			return nil
		},
	}
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

	if err := service.CreateRecord(title, record); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(cs.Green("Record has been successfully created with ID %q", record.ID))

	for _, id := range opts.supersedes {
		rcd, ok := service.GetRecord(id)
		if !ok {
			continue
		}
		rcd.Status = records.SUPERSEDED
		rcd.Superseders.Append(record.ID)
		if err := service.UpdateRecord(rcd); err != nil {
			printWarning("Unable to update record %q: %v", rcd.ID, err)
		}
	}

	fmt.Println()
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
