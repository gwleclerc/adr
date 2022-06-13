package records

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gernest/front"
	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/ojizero/gofindup"
	"golang.org/x/sync/semaphore"
	"gopkg.in/yaml.v3"
)

var matter = front.NewMatter()

func init() {
	matter.Handle("---", front.YAMLHandler)
}

func retrieveADRsPath() (string, error) {
	path, err := gofindup.Findup(cs.ConfigurationFile)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	var config cs.Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(filepath.Dir(path), config.Directory)
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%q should be a directory", fullPath)
	}
	return fullPath, nil
}

func indexADRs(path string) ([]AdrData, error) {
	res := []AdrData{}

	mu := sync.Mutex{}
	sem := semaphore.NewWeighted(10)
	wg := sync.WaitGroup{}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q should be a directory", path)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	for _, f := range files {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			fmt.Println(cs.Yellow("Unable to acquire semaphore: %s", err.Error()))
			continue
		}
		if f.IsDir() {
			continue
		}
		wg.Add(1)
		go func(f fs.FileInfo) {
			defer wg.Done()
			defer sem.Release(1)

			adrData := AdrData{
				Name: f.Name(),
			}

			filePath := filepath.Join(path, adrData.Name)
			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println(cs.Yellow("Unable to read file %q: %s", filePath, err.Error()))
				return
			}

			data, body, err := matter.Parse(bytes.NewReader(b))
			if err != nil {
				fmt.Println(cs.Yellow("Unable to read yaml header from file %q: %s", filePath, err.Error()))
				return
			}
			adrData.Body = body

			err = processDate(data, "creation_date", adrData.Name)
			if err != nil {
				fmt.Println(cs.Yellow("Invalid creation date in yaml header from file %q: %v", filePath, err))
				return
			}
			err = processDate(data, "last_update_date", adrData.Name)
			if err != nil {
				fmt.Println(cs.Yellow("Invalid last update date in yaml header from file %q: %v", filePath, err))
				return
			}
			processSet(data, "tags")
			processSet(data, "superseders")

			err = mapstructure.Decode(data, &adrData)
			if err != nil {
				fmt.Println(cs.Yellow("Invalid yaml header: %s", filePath, err.Error()))
				return
			}

			mu.Lock()
			defer mu.Unlock()
			res = append(res, adrData)
		}(f)
	}
	wg.Wait()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res, nil
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
