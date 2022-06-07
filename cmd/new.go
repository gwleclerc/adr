package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	simpleSlug "github.com/gosimple/slug"
	. "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/templates"
	"github.com/gwleclerc/adr/types"
	"github.com/gwleclerc/adr/utils"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"github.com/teris-io/shortid"
)

var (
	new_author string
	new_tags   []string
	new_status types.AdrStatus = types.ACCEPTED
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [flags] <record title...>",
	Short: "Create a new ADR",
	Long: fmt.Sprintf(
		`
Create a new architecture decision record.
It will be created in the directory defined in the nearest %s configuration file.`,
		ConfigurationFile,
	),
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.TrimSpace(strings.Join(args, " "))
		if title == "" {
			fmt.Printf("%s %s %s\n", Red("invalid argument: please specify a"), RedUnderline("title"), Red("as arguments"))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		path, err := utils.RetrieveADRsPath()
		if err != nil {
			fmt.Println(Red("unable to retrieve ADRs path, you should look at the %s configuration file: %v", ConfigurationFile, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		if err := newRecord(path, title); err != nil {
			fmt.Println(Red("unable to create a new ADRs in directory %q: %v", path, err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
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
	newCmd.RegisterFlagCompletionFunc("status", types.AdrStatusCompletion)
	newCmd.Flags().StringSliceVarP(
		&new_tags,
		"tags",
		"t",
		[]string{},
		`tags of the record`,
	)
	rootCmd.AddCommand(newCmd)
}

func newRecord(path, title string) error {
	adrs, err := utils.IndexADRs(path)
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("%03d", 1)
	for i := range adrs {
		record := adrs[len(adrs)-1-i]
		if number := utils.GetRecordNumber(record.Name); number != "" {
			count, _ := strconv.Atoi(number)
			prefix = fmt.Sprintf("%03d", count+1)
			break
		}
	}

	slug := strings.ReplaceAll(simpleSlug.Make(title), "-", "_")
	filename := fmt.Sprintf("%s_%s.md", prefix, slug)

	filePath := filepath.Join(path, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if new_author == "" {
		username, err := gitconfig.Username()
		if err != nil {
			fmt.Println(Yellow("Unable to find a git user: %v", err))
			user, err := user.Current()
			if err != nil {
				fmt.Printf(Yellow("Unable to find a OS user: %v", err))
				username = DefaultUserName
			} else {
				username = user.Username
			}
		}
		new_author = username
	}

	date := time.Now()
	record := types.AdrData{
		ID:             shortid.MustGenerate(),
		Title:          slug,
		Status:         new_status,
		CreationDate:   date,
		LastUpdateDate: date,
		Author:         new_author,
		Tags:           make(types.Set[string]),
	}
	record.Tags.Append(new_tags...)

	b, err := types.MarshalYAML(record)
	if err != nil {
		return err
	}

	humanizedCreationDate := record.CreationDate.Format(time.RFC1123)
	err = templates.Templates[CreateADRTemplate].Execute(file, map[string]any{
		"Header": strings.Trim(string(b), "\n"),
		"Title":  strings.Title(strings.ToLower(title)),
		"Date":   humanizedCreationDate,
	})
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println(Green("Record has been successfully created at %q with ID %q", filePath, record.ID))
	fmt.Println()
	return nil
}
