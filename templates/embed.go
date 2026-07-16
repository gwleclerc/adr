package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed record.tpl bodies/*.tpl
var files embed.FS

// recordTmpl is the shared record envelope: YAML front-matter + body.
var recordTmpl = template.Must(template.New("record").Parse(mustRead("record.tpl")))

func mustRead(name string) string {
	b, err := files.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Template is a named ADR body skeleton (the section headings + their guidance).
type Template struct {
	Name    string
	Body    string
	Builtin bool
}

// RenderRecord wraps a body with the record envelope (front-matter + body).
func RenderRecord(header, body string) (string, error) {
	var sb strings.Builder
	err := recordTmpl.Execute(&sb, map[string]any{
		"Header": strings.Trim(header, "\n"),
		"Body":   strings.TrimRight(body, "\n"),
	})
	return sb.String(), err
}

// Builtins returns the templates embedded in the binary, keyed by name.
func Builtins() map[string]Template {
	out := map[string]Template{}
	entries, err := fs.ReadDir(files, "bodies")
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tpl") {
			continue
		}
		name := templateName(e.Name())
		out[name] = Template{Name: name, Body: mustRead("bodies/" + e.Name()), Builtin: true}
	}
	return out
}

// Load returns the built-in templates merged with any custom `*.tpl` found in
// customDir (whose names override built-ins). A missing customDir is ignored.
func Load(customDir string) (map[string]Template, error) {
	out := Builtins()
	if customDir == "" {
		return out, nil
	}
	entries, err := os.ReadDir(customDir)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tpl") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(customDir, e.Name()))
		if err != nil {
			return nil, err
		}
		name := templateName(e.Name())
		out[name] = Template{Name: name, Body: string(b), Builtin: false}
	}
	return out, nil
}

// Names returns the template names sorted alphabetically.
func Names(templates map[string]Template) []string {
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Headings extracts the markdown heading lines (## and deeper) from a body.
func Headings(body string) []string {
	var hs []string
	for _, line := range strings.Split(body, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "#") {
			hs = append(hs, t)
		}
	}
	return hs
}

// Validate checks that providedBody contains every heading of the template body,
// in the same order, each followed by some non-blank content.
func Validate(templateBody, providedBody string) error {
	want := Headings(templateBody)
	lines := strings.Split(providedBody, "\n")

	// Locate each wanted heading in order.
	positions := make([]int, 0, len(want))
	idx := 0
	for i, line := range lines {
		if idx < len(want) && strings.TrimSpace(line) == want[idx] {
			positions = append(positions, i)
			idx++
		}
	}
	if idx < len(want) {
		return fmt.Errorf("missing or out-of-order section %q", strings.TrimSpace(want[idx]))
	}

	// Each section must have at least one non-blank content line before the next.
	for k, start := range positions {
		end := len(lines)
		if k+1 < len(positions) {
			end = positions[k+1]
		}
		hasContent := false
		for _, line := range lines[start+1 : end] {
			if strings.TrimSpace(line) != "" {
				hasContent = true
				break
			}
		}
		if !hasContent {
			return fmt.Errorf("section %q is empty", strings.TrimSpace(want[k]))
		}
	}
	return nil
}

func templateName(filename string) string {
	return strings.TrimSuffix(strings.ToLower(filepath.Base(filename)), ".tpl")
}
