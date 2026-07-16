package records

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/gernest/front"
	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/ojizero/gofindup"
	"gopkg.in/yaml.v3"
)

var matter = front.NewMatter()

func init() {
	matter.Handle("---", front.YAMLHandler)
}

// LoadConfig finds the nearest configuration file and returns the parsed config
// along with the directory that contains it (used to resolve relative paths).
func LoadConfig() (cs.Config, string, error) {
	path, err := gofindup.Findup(cs.ConfigurationFile)
	if err != nil {
		return cs.Config{}, "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return cs.Config{}, "", err
	}
	var config cs.Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		return cs.Config{}, "", err
	}
	return config, filepath.Dir(path), nil
}

func indexADRs(path string) ([]AdrData, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	res := make([]AdrData, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if adr, ok := parseADR(path, entry.Name()); ok {
			res = append(res, adr)
		}
	}
	slices.SortFunc(res, func(a, b AdrData) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return res, nil
}

// parseADR reads and parses a single ADR file. It returns ok=false (after logging
// a warning to stderr) when the file cannot be read or parsed, so one bad file
// does not abort indexing.
func parseADR(dir, name string) (AdrData, bool) {
	adrData := AdrData{Name: name}
	filePath := filepath.Join(dir, name)

	b, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, cs.Yellow("Unable to read file %q: %s", filePath, err.Error()))
		return AdrData{}, false
	}

	data, body, err := matter.Parse(bytes.NewReader(b))
	if err != nil {
		fmt.Fprintln(os.Stderr, cs.Yellow("Unable to read yaml header from file %q: %s", filePath, err.Error()))
		return AdrData{}, false
	}
	adrData.Body = body

	if err := processDate(data, "creation_date", name); err != nil {
		fmt.Fprintln(os.Stderr, cs.Yellow("Invalid creation date in yaml header from file %q: %v", filePath, err))
		return AdrData{}, false
	}
	if err := processDate(data, "last_update_date", name); err != nil {
		fmt.Fprintln(os.Stderr, cs.Yellow("Invalid last update date in yaml header from file %q: %v", filePath, err))
		return AdrData{}, false
	}
	processSet(data, "tags")
	processSet(data, "superseders")

	if err := mapstructure.Decode(data, &adrData); err != nil {
		fmt.Fprintln(os.Stderr, cs.Yellow("Invalid yaml header in file %q: %v", filePath, err))
		return AdrData{}, false
	}
	return adrData, true
}

func processDate(data map[string]any, dateKey, recordName string) error {
	// If the date is missing, we init it with the zero value "0001-01-01 00:00:00 +0000 UTC"
	if data[dateKey] == nil {
		data[dateKey] = new(time.Time)
	}
	dateTime, ok := data[dateKey].(*time.Time)
	if ok {
		// If data[dateKey] is a *time.Time it is necessarily the zero value
		// so we arbitrary add the number of the record as days to keep records order
		if number := utils.GetRecordNumber(recordName); number != "" {
			num, _ := strconv.Atoi(number)
			data[dateKey] = dateTime.Add(time.Duration(num) * 24 * time.Hour)
		}
		return nil
	}
	date, ok := data[dateKey].(string)
	if !ok {
		return errors.New("invalid date value")
	}

	var err error
	data[dateKey], err = time.Parse(time.RFC3339, date)
	if err != nil {
		return err
	}
	return nil
}

func processSet(data map[string]any, key string) {
	unknown, ok := data[key].([]any)
	if !ok {
		return
	}

	tmp := make([]string, 0, len(unknown))
	for _, elem := range unknown {
		tmp = append(tmp, fmt.Sprintf("%v", elem))
	}
	set := make(Set[string], len(unknown))
	set.Append(tmp...)
	data[key] = set
}
