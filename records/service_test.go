package records

import (
	"os"
	"path/filepath"
	"testing"
)

// newTestService sets up an isolated ADR project in a temp dir and returns a service.
func newTestProject(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "adrs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".adrrc.yml"), []byte("directory: adrs\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
}

func TestServiceLifecycle(t *testing.T) {
	newTestProject(t)

	svc, err := NewService()
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}

	// Create numbers and slugs the file, and keeps the human title verbatim.
	rec := AdrData{ID: "aaa", Status: ACCEPTED, Author: "me", Tags: make(Set[string])}
	created, err := svc.CreateRecord("My First Decision", rec, "## Context\nbecause\n")
	if err != nil {
		t.Fatalf("CreateRecord: %v", err)
	}
	if created.Name != "001_my_first_decision.md" {
		t.Errorf("filename = %q, want 001_my_first_decision.md", created.Name)
	}
	if created.Title != "My First Decision" {
		t.Errorf("title = %q, want verbatim", created.Title)
	}

	// A fresh service re-indexes the record from disk.
	svc2, err := NewService()
	if err != nil {
		t.Fatalf("NewService (reindex): %v", err)
	}
	got := svc2.GetRecords()
	if len(got) != 1 || got[0].ID != "aaa" {
		t.Fatalf("GetRecords = %+v, want one record with ID aaa", got)
	}

	// The next record gets the incremented numeric prefix.
	second, err := svc2.CreateRecord("Second Decision", AdrData{ID: "bbb", Status: ACCEPTED, Tags: make(Set[string])}, "## Context\nmore\n")
	if err != nil {
		t.Fatalf("CreateRecord second: %v", err)
	}
	if second.Name != "002_second_decision.md" {
		t.Errorf("second filename = %q, want 002_second_decision.md", second.Name)
	}

	// Update persists a metadata change.
	svc3, _ := NewService()
	r, ok := svc3.GetRecord("aaa")
	if !ok {
		t.Fatal("record aaa not found")
	}
	r.Status = DEPRECATED
	if err := svc3.UpdateRecord(r); err != nil {
		t.Fatalf("UpdateRecord: %v", err)
	}
	svc4, _ := NewService()
	if r4, _ := svc4.GetRecord("aaa"); r4.Status != DEPRECATED {
		t.Errorf("status after update = %q, want deprecated", r4.Status)
	}
}

func TestServiceIgnoresNonRecords(t *testing.T) {
	newTestProject(t)

	svc, err := NewService()
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	if _, err := svc.CreateRecord("A decision", AdrData{ID: "a", Status: ACCEPTED, Tags: make(Set[string])}, "## Context\nx\n"); err != nil {
		t.Fatalf("CreateRecord: %v", err)
	}

	// Drop non-record files in the ADR directory (e.g. a generated index).
	adrs := filepath.Join(".", "adrs")
	if err := os.WriteFile(filepath.Join(adrs, "README.md"), []byte("# Index\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adrs, "notes.txt"), []byte("scratch"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc2, err := NewService()
	if err != nil {
		t.Fatalf("NewService (reindex): %v", err)
	}
	if got := svc2.GetRecords(); len(got) != 1 {
		t.Errorf("expected 1 record (non-records ignored), got %d: %+v", len(got), got)
	}
}
