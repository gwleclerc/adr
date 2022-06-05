/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
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
	"github.com/gwleclerc/adr/utils"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"github.com/teris-io/shortid"
	"gopkg.in/yaml.v3"
)

var (
	author string
	tags   []string
	status AdrStatus = ACCEPTED
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [flags] <record title...>",
	Short: "Add a new ADR",
	Long: fmt.Sprintf(
		`
Create a new architecture decision record.
It will be created in the directory defined in the nearest %s configuration file.`,
		ConfigurationFile,
	),
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args, " ")
		if len(args) <= 0 || title == "" {
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
		fmt.Println(Green("Creating a new record %q", title))
		if err := addRecord(path, title); err != nil {
			fmt.Println(Red("unable to add ADRs directory: %v", err))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		cmd.Println(Green("Record has been successfully created at %q", path))
	},
}

func init() {
	addCmd.Flags().StringVarP(
		&author,
		"author",
		"a",
		"",
		"author of the record",
	)
	addCmd.Flags().VarP(
		&status,
		"status",
		"s",
		`status of the record, allowed: "unknown", "proposed", "accepted", "deprecated" or "superseded"`,
	)
	addCmd.Flags().StringSliceVarP(
		&tags,
		"tags",
		"t",
		[]string{},
		`tags of the record`,
	)

	rootCmd.AddCommand(addCmd)
}

func addRecord(path, title string) error {
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
	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	if author == "" {
		username, err := gitconfig.Username()
		if err != nil {
			fmt.Printf("Unable to find a git user:\n\t%s", err.Error())
			user, err := user.Current()
			if err != nil {
				fmt.Printf("Unable to find a OS user:\n\t%s", err.Error())
				username = DefaultUserName
			} else {
				username = user.Username
			}
		}
		author = username
	}

	record := AdrData{
		ID:     shortid.MustGenerate(),
		Title:  slug,
		Status: status,
		Date:   time.Now(),
		Author: author,
		Tags:   tags,
	}

	b, err := yaml.Marshal(record)
	if err != nil {
		return err
	}

	humanizedDate := record.Date.Format(time.RFC1123)
	err = templates.Templates[CreateADRTemplate].Execute(file, map[string]any{
		"Header": strings.Trim(string(b), "\n"),
		"Title":  strings.Title(strings.ToLower(title)),
		"Date":   humanizedDate,
	})
	if err != nil {
		return err
	}

	return nil
}
