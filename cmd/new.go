package cmd

import (
	"fmt"
	"os/user"
	"strings"

	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/records"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"github.com/teris-io/shortid"
)

var (
	new_author     string
	new_tags       []string
	new_status     records.AdrStatus = records.ACCEPTED
	new_supersedes []string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [flags] <record title...>",
	Short: "Create a new ADR",
	Long: fmt.Sprintf(
		`
Create a new architecture decision record.
It will be created in the directory defined in the nearest %s configuration file.`,
		cs.ConfigurationFile,
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.TrimSpace(strings.Join(args, " "))
		if title == "" {
			fmt.Printf("%s %s %s\n", cs.Red("invalid argument: please specify a"), cs.RedUnderline("title"), cs.Red("in arguments"))
			return ErrSilent
		}
		service, err := records.NewService()
		if err != nil {
			fmt.Println(cs.Red("unable to initialize records service: %v", err))
			return ErrSilent
		}
		if err := newRecord(service, title); err != nil {
			fmt.Println(cs.Red("unable to create a new ADRs: %v", err))
			return ErrSilent
		}
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(
		&new_author,
		"author",
		"a",
		"",
		"author of the record",
	)
	newCmd.Flags().VarP(
		&new_status,
		"status",
		"s",
		`status of the record, allowed: "unknown", "proposed", "accepted", "deprecated" or "superseded"`,
	)
	_ = newCmd.RegisterFlagCompletionFunc("status", records.AdrStatusCompletion)
	newCmd.Flags().StringSliceVarP(
		&new_tags,
		"tags",
		"t",
		[]string{},
		`tags of the record`,
	)
	newCmd.Flags().StringSliceVarP(
		&new_supersedes,
		"supersedes",
		"r",
		[]string{},
		`record ids superseded by this one`,
	)
	rootCmd.AddCommand(newCmd)
}

func newRecord(service *records.Service, title string) error {
	if new_author == "" {
		username, err := gitconfig.Username()
		if err != nil {
			fmt.Println(cs.Yellow("Unable to find a git user: %v", err))
			user, err := user.Current()
			if err != nil {
				fmt.Println(cs.Yellow("Unable to find a OS user: %v", err))
				username = cs.DefaultUserName
			} else {
				username = user.Username
			}
		}
		new_author = username
	}

	// Since IDs starting with '-' will be interpreted as CLI flags, we have to regenerate a new ID until this is no longer the case.
	id := shortid.MustGenerate()
	for strings.HasPrefix(id, "-") {
		id = shortid.MustGenerate()
	}

	record := records.AdrData{
		ID:     id,
		Status: new_status,
		Author: new_author,
		Tags:   make(records.Set[string]),
	}
	record.Tags.Append(new_tags...)

	err := service.CreateRecord(title, record)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(cs.Green("Record has been successfully created with ID %q", record.ID))

	for _, id := range new_supersedes {
		rcd, ok := service.GetRecord(id)
		if !ok {
			continue
		}
		rcd.Status = records.SUPERSEDED
		rcd.Superseders.Append(record.ID)
		if err := service.UpdateRecord(rcd); err != nil {
			fmt.Println(cs.Yellow("Unable to update record %q: %v", rcd.ID, err))
		}
	}

	fmt.Println()
	return nil
}
