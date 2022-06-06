package utils

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
	. "github.com/gwleclerc/adr/constants"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"github.com/ojizero/gofindup"
	"golang.org/x/sync/semaphore"
	"gopkg.in/yaml.v3"
)

var matter = front.NewMatter()

func init() {
	matter.Handle("---", front.YAMLHandler)
}

func RetrieveADRsPath() (string, error) {
	path, err := gofindup.Findup(ConfigurationFile)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	var config Config
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

func IndexADRs(path string) ([]AdrData, error) {
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
			fmt.Println(Yellow("Unable to acquire semaphore: %s", err.Error()))
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
				fmt.Println(Yellow("Unable to read file %q: %s", filePath, err.Error()))
				return
			}

			data, body, err := matter.Parse(bytes.NewReader(b))
			if err != nil {
				fmt.Println(Yellow("Unable to read yaml header from file %q: %s", filePath, err.Error()))
				return
			}
			adrData.Body = body

			err = processDate(data, "creation_date", adrData.Name)
			if err != nil {
				fmt.Println(Yellow("Invalid creation date in yaml header from file %q: %v", filePath, err))
				return
			}
			err = processDate(data, "last_update_date", adrData.Name)
			if err != nil {
				fmt.Println(Yellow("Invalid last update date in yaml header from file %q: %v", filePath, err))
				return
			}
			if err != nil {
				fmt.Println(Yellow("Unable to parse date from yaml header in file %q: %s", filePath, err.Error()))
				return
			}
			mapstructure.Decode(data, &adrData)
			if err != nil {
				fmt.Println(Yellow("Invalid yaml header: %s", filePath, err.Error()))
				return
			}

			mu.Lock()
			defer mu.Unlock()
			res = append(res, adrData)
		}(f)
	}
	wg.Wait()
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreationDate.Before(res[j].CreationDate)
	})
	return res, nil
}

func processDate(data map[string]any, dateKey, recordName string) error {
	// We must convert the key to camel case because mapstructure does not process snake case keys correctly
	destKey := strcase.ToCamel(dateKey)
	// If the date is missing, we init it with the zero value "0001-01-01 00:00:00 +0000 UTC"
	if data[dateKey] == nil {
		data[dateKey] = new(time.Time)
	}
	dateTime, ok := data[dateKey].(*time.Time)
	if ok {
		// If data[dateKey] is a *time.Time it is necessarily the zero value
		// so we arbitrary add the number of the record as days to keep records order
		if number := GetRecordNumber(recordName); number != "" {
			num, _ := strconv.Atoi(number)
			data[destKey] = dateTime.Add(time.Duration(num) * 24 * time.Hour)
		}
		return nil
	}
	date, ok := data[dateKey].(string)
	if !ok {
		return errors.New("invalid date value")
	}

	var err error
	data[destKey], err = time.Parse(time.RFC3339, date)
	if err != nil {
		return err
	}
	return nil
}
