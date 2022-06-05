package templates

import (
	"embed"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed *
var templates embed.FS

var Templates = map[string]*template.Template{}

func init() {
	err := fs.WalkDir(templates, ".", func(path string, info fs.DirEntry, err error) error {
		// Skip non-templates.
		if info.IsDir() || !strings.HasSuffix(path, ".tpl") {
			return nil
		}
		name := strings.ToLower(filepath.Base(path))
		tpl := template.New(name)
		// Load file from embed virtual file, or use the shortcut
		// templates.ReadFile(path).
		f, err := templates.Open(path)
		if err != nil {
			return err
		}
		// Now read it.
		sl, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		// It can now be parsed as a string.
		tpl, err = tpl.Parse(string(sl))
		if err != nil {
			return err
		}
		Templates[name] = tpl
		return nil
	})
	if err != nil {
		panic(err)
	}
}
