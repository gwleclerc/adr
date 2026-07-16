package templates

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuiltins(t *testing.T) {
	b := Builtins()
	for _, name := range []string{"bare", "madr"} {
		tpl, ok := b[name]
		if !ok {
			t.Errorf("missing built-in template %q", name)
			continue
		}
		if !tpl.Builtin {
			t.Errorf("template %q should be flagged as built-in", name)
		}
	}
}

func TestValidate(t *testing.T) {
	madr := Builtins()["madr"].Body

	valid := "## Context and Problem Statement\nx\n\n## Considered Options\ny\n\n## Decision Outcome\nz\n\n### Consequences\nw\n"
	if err := Validate(madr, valid); err != nil {
		t.Errorf("valid body rejected: %v", err)
	}

	missing := "## Context and Problem Statement\nx\n"
	if err := Validate(madr, missing); err == nil {
		t.Error("body missing sections should be rejected")
	}

	empty := "## Context and Problem Statement\n\n## Considered Options\ny\n\n## Decision Outcome\nz\n\n### Consequences\nw\n"
	if err := Validate(madr, empty); err == nil {
		t.Error("body with an empty section should be rejected")
	}

	reordered := "## Considered Options\ny\n\n## Context and Problem Statement\nx\n\n## Decision Outcome\nz\n\n### Consequences\nw\n"
	if err := Validate(madr, reordered); err == nil {
		t.Error("out-of-order body should be rejected")
	}
}

func TestLoad(t *testing.T) {
	// No custom dir: built-ins only.
	reg, err := Load("")
	if err != nil {
		t.Fatalf("Load(\"\") error: %v", err)
	}
	if _, ok := reg["bare"]; !ok {
		t.Error("built-ins should be present with no custom dir")
	}

	// Custom dir: a new template and an override of a built-in.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "custom.tpl"), []byte("## X\n> guidance\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "madr.tpl"), []byte("## Overridden\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg, err = Load(dir)
	if err != nil {
		t.Fatalf("Load(%q) error: %v", dir, err)
	}
	if _, ok := reg["custom"]; !ok {
		t.Error("custom template should be loaded")
	}
	if reg["custom"].Builtin {
		t.Error("custom template should not be flagged built-in")
	}
	if reg["madr"].Body != "## Overridden\n" {
		t.Errorf("madr should be overridden by the custom dir, got %q", reg["madr"].Body)
	}
}
