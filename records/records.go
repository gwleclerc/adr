package records

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	simpleSlug "github.com/gosimple/slug"
	cs "github.com/gwleclerc/adr/constants"
	"github.com/gwleclerc/adr/templates"
	"github.com/gwleclerc/adr/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service struct {
	records  map[string]AdrData
	ids      []string
	adrsPath string
}

func NewService() (*Service, error) {
	path, err := retrieveADRsPath()
	if err != nil {
		return nil, err
	}
	adrs, err := indexADRs(path)
	if err != nil {
		return nil, err
	}
	records := make(map[string]AdrData, len(adrs))
	ids := make([]string, 0, len(adrs))
	for _, adr := range adrs {
		records[adr.ID] = adr
		ids = append(ids, adr.ID)
	}
	return &Service{
		records:  records,
		ids:      ids,
		adrsPath: path,
	}, nil
}

func (s Service) GetRecord(recordID string) (AdrData, bool) {
	adr, ok := s.records[recordID]
	return adr, ok
}

func (s Service) GetRecords() []AdrData {
	records := make([]AdrData, 0, len(s.records))
	for _, id := range s.ids {
		records = append(records, s.records[id])
	}
	return records
}

func (s Service) CreateRecord(title string, record AdrData) error {
	prefix := fmt.Sprintf("%03d", 1)
	for i := range s.ids {
		recordID := s.ids[len(s.ids)-1-i]
		record := s.records[recordID]
		if number := utils.GetRecordNumber(record.Name); number != "" {
			count, _ := strconv.Atoi(number)
			prefix = fmt.Sprintf("%03d", count+1)
			break
		}
	}

	slug := strings.ReplaceAll(simpleSlug.Make(title), "-", "_")
	filename := fmt.Sprintf("%s_%s.md", prefix, slug)

	date := time.Now()
	record.Title = slug
	record.CreationDate = date
	record.LastUpdateDate = date

	b, err := MarshalYAML(record)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(s.adrsPath, filename))
	if err != nil {
		return err
	}
	defer file.Close()
	err = templates.Templates[cs.CreateADRTemplate].Execute(file, map[string]any{
		"Header": strings.Trim(string(b), "\n"),
		"Title":  cases.Title(language.Und).String(strings.ToLower(title)),
		"Date":   record.CreationDate.Format(time.RFC1123),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s Service) UpdateRecord(record AdrData) error {
	record.LastUpdateDate = time.Now()

	b, err := MarshalYAML(record)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(s.adrsPath, record.Name))
	if err != nil {
		return err
	}
	defer file.Close()

	err = templates.Templates[cs.UpdateADRTemplate].Execute(file, map[string]any{
		"Header": strings.Trim(string(b), "\n"),
		"Body":   record.Body,
	})
	if err != nil {
		return err
	}

	return nil
}
