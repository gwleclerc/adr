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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service struct {
	records      map[string]AdrData
	ids          []string
	adrsPath     string
	templatesDir string
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
		records:      records,
		ids:          ids,
		adrsPath:     adrsPath,
		templatesDir: templatesDir,
	}, nil
}

// TemplatesDir returns the resolved custom templates directory ("" if unset).
func (s Service) TemplatesDir() string {
	return s.templatesDir
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
func (s Service) CreateRecord(title string, record AdrData, body string) error {
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

	slug := strings.ReplaceAll(simpleSlug.Make(title), "-", "_")
	filename := fmt.Sprintf("%s_%s.md", prefix, slug)

	date := time.Now()
	record.Title = slug
	record.CreationDate = date
	record.LastUpdateDate = date

	header, err := MarshalYAML(record)
	if err != nil {
		return err
	}

	titleCased := cases.Title(language.Und).String(strings.ToLower(title))
	fullBody := fmt.Sprintf("# %s\n\nDate: %s\n\n%s",
		titleCased, date.Format(time.RFC1123), strings.TrimRight(body, "\n"))

	return s.writeRecord(filename, string(header), fullBody)
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
	file, err := os.Create(filepath.Join(s.adrsPath, filename))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(out)
	return err
}
