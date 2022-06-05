package utils

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/gernest/front"
	. "github.com/gwleclerc/adr/constants"
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
			filePath := filepath.Join(path, f.Name())
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
			date, ok := data["date"].(string)
			if !ok {
				fmt.Println(Yellow("Invalid date in yaml header from file %q", filePath))
				return
			}
			data["date"], err = time.Parse(time.RFC3339, date)
			if err != nil {
				fmt.Println(Yellow("Unable to parse date from yaml header in file %q: %s", filePath, err.Error()))
				return
			}
			var adrData AdrData
			mapstructure.Decode(data, &adrData)
			if err != nil {
				fmt.Println(Yellow("Invalid yaml header: %s", filePath, err.Error()))
				return
			}
			adrData.Name = f.Name()
			adrData.Body = body

			mu.Lock()
			defer mu.Unlock()
			res = append(res, adrData)
		}(f)
	}
	wg.Wait()
	mu.Lock()
	defer mu.Unlock()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Date.Before(res[j].Date)
	})
	return res, nil
}
