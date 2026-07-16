package records

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	simpleSlug "github.com/gosimple/slug"
	"github.com/gwleclerc/adr/templates"
	"github.com/gwleclerc/adr/utils"
)

type Service struct {
	records         map[string]AdrData
	ids             []string
	adrsPath        string
	templatesDir    string
	defaultTemplate string
	defaultAuthor   string
}

func NewService() (*Service, error) {
	cfg, dir, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	adrsPath := filepath.Join(dir, cfg.Directory)
	info, err := os.Stat(adrsPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q should be a directory", adrsPath)
	}

	adrs, err := indexADRs(adrsPath)
	if err != nil {
		return nil, err
	}
	records := make(map[string]AdrData, len(adrs))
	ids := make([]string, 0, len(adrs))
	for _, adr := range adrs {
		records[adr.ID] = adr
		ids = append(ids, adr.ID)
	}

	templatesDir := ""
	if cfg.TemplatesDir != "" {
		templatesDir = filepath.Join(dir, cfg.TemplatesDir)
	}
	return &Service{
		records:         records,
		ids:             ids,
		adrsPath:        adrsPath,
		templatesDir:    templatesDir,
		defaultTemplate: cfg.DefaultTemplate,
		defaultAuthor:   cfg.DefaultAuthor,
	}, nil
}

// TemplatesDir returns the resolved custom templates directory ("" if unset).
func (s Service) TemplatesDir() string {
	return s.templatesDir
}

// DefaultTemplate returns the template configured as default ("" if unset).
func (s Service) DefaultTemplate() string {
	return s.defaultTemplate
}

// DefaultAuthor returns the author configured as default ("" if unset).
func (s Service) DefaultAuthor() string {
	return s.defaultAuthor
}

// RecordPath returns the absolute path of a record's file.
func (s Service) RecordPath(record AdrData) string {
	return filepath.Join(s.adrsPath, record.Name)
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

// CreateRecord writes a new record file. body is the markdown of the sections
// (a template skeleton or a caller-provided, already-validated body); the record
// envelope (front-matter, title and date) is added here.
func (s Service) CreateRecord(title string, record AdrData, body string) (AdrData, error) {
	prefix := fmt.Sprintf("%03d", 1)
	for i := range s.ids {
		recordID := s.ids[len(s.ids)-1-i]
		previous := s.records[recordID]
		if number := utils.GetRecordNumber(previous.Name); number != "" {
			count, _ := strconv.Atoi(number)
			prefix = fmt.Sprintf("%03d", count+1)
			break
		}
	}

	title = strings.TrimSpace(title)
	slug := strings.ReplaceAll(simpleSlug.Make(title), "-", "_")
	filename := fmt.Sprintf("%s_%s.md", prefix, slug)

	date := time.Now()
	// Store the human-readable title in the metadata; the slug lives only in the
	// filename. The title is used verbatim so acronyms and casing are preserved.
	record.Title = title
	record.CreationDate = date
	record.LastUpdateDate = date
	record.Name = filename

	header, err := MarshalYAML(record)
	if err != nil {
		return AdrData{}, err
	}

	fullBody := fmt.Sprintf("# %s\n\nDate: %s\n\n%s",
		title, date.Format(time.RFC1123), strings.TrimRight(body, "\n"))

	if err := s.writeRecord(filename, string(header), fullBody); err != nil {
		return AdrData{}, err
	}
	return record, nil
}

func (s Service) UpdateRecord(record AdrData) error {
	record.LastUpdateDate = time.Now()

	header, err := MarshalYAML(record)
	if err != nil {
		return err
	}

	return s.writeRecord(record.Name, string(header), record.Body)
}

func (s Service) writeRecord(filename, header, body string) error {
	out, err := templates.RenderRecord(header, body)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.adrsPath, filename), []byte(out), 0o644)
}
