package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anmitsu/goful/filer"
)

func resetExcludeState(t *testing.T) {
	t.Helper()
	filer.ConfigureExcludedNames(nil, false)
}

func TestParseExcludeLines_REQ_FILER_EXCLUDE_NAMES(t *testing.T) {
	// [REQ:FILER_EXCLUDE_NAMES] [ARCH:FILER_EXCLUDE_FILTER] [IMPL:FILER_EXCLUDE_LOADER]
	input := `
# comment

.DS_Store
Thumbs.DB

`
	names, err := parseExcludeLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseExcludeLines error: %v", err)
	}
	if len(names) != 2 || names[0] != ".DS_Store" || names[1] != "Thumbs.DB" {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestLoadExcludedNames_REQ_FILER_EXCLUDE_NAMES(t *testing.T) {
	// [REQ:FILER_EXCLUDE_NAMES] [ARCH:FILER_EXCLUDE_FILTER] [IMPL:FILER_EXCLUDE_LOADER]
	t.Run("missing file disables filter", func(t *testing.T) {
		resetExcludeState(t)
		loadExcludedNames(filepath.Join(t.TempDir(), "missing"))
		if filer.ExcludedNamesEnabled() {
			t.Fatalf("filter should remain disabled when file is missing")
		}
	})

	t.Run("valid file enables filter", func(t *testing.T) {
		resetExcludeState(t)
		dir := t.TempDir()
		path := filepath.Join(dir, "excludes")
		if err := os.WriteFile(path, []byte(".DS_Store\nThumbs.db\n"), 0o644); err != nil {
			t.Fatalf("write exclude file: %v", err)
		}
		loadExcludedNames(path)
		if !filer.ExcludedNamesEnabled() {
			t.Fatalf("filter should be enabled after loading entries")
		}
	})
}
