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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <record title...>",
	Short: "Add a new ADR",
	Long: `Create a new architecture decision record.
	It will be created in the directory defined in the nearest .adrrc configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args, " ")
		if len(args) <= 0 || title == "" {
			fmt.Printf("%s %s %s\n", Red("You must specify a"), RedUnderline("title"), Red("as arguments."))
			cmd.Usage()
			os.Exit(1)
		}
		path, err := utils.RetrieveADRsPath()
		if err != nil {
			fmt.Println(Red("Unable to retrieve ADRs path, you should look at the .adrrc configuration file:\n\t%s", err.Error()))
			os.Exit(1)
		}
		fmt.Println(Green("Creating a new record %q", title))
		if err := addRecord(path, title); err != nil {
			fmt.Println(Red("Unable to add ADRs directory:\n\t%s", err.Error()))
		}
		fmt.Println(Green("Record has been successfully created at %q", path))
	},
}

func init() {
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

	record := AdrData{
		ID:     shortid.MustGenerate(),
		Title:  slug,
		Status: ACCEPTED,
		Date:   time.Now(),
		Author: username,
	}

	b, err := yaml.Marshal(record)
	if err != nil {
		return err
	}

	humanizedDate := record.Date.Format(time.RFC1123)
	err = templates.Templates[CreateADRTemplate].Execute(file, map[string]any{
		"Header": strings.Trim(string(b), "\n"),
		"Title":  title,
		"Date":   humanizedDate,
	})
	if err != nil {
		return err
	}

	return nil
}
