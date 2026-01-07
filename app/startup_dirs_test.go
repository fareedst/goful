package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseStartupDirs_REQ_WORKSPACE_START_DIRS(t *testing.T) {
	t.Parallel()

	temp := t.TempDir()
	dirA := filepath.Join(temp, "dirA")
	dirB := filepath.Join(temp, "dirB")
	if err := os.MkdirAll(dirA, 0o755); err != nil {
		t.Fatalf("mkdir dirA: %v", err)
	}
	if err := os.MkdirAll(dirB, 0o755); err != nil {
		t.Fatalf("mkdir dirB: %v", err)
	}

	filePath := filepath.Join(temp, "file.txt")
	if err := os.WriteFile(filePath, []byte("noop"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	missing := filepath.Join(temp, "missing")
	dirs, warnings := ParseStartupDirs([]string{"  ", dirA, missing, dirB, filePath, dirA, dirB})

	if len(dirs) != 4 {
		t.Fatalf("expected 4 directories (duplicates allowed), got %d (%v)", len(dirs), dirs)
	}
	if dirs[0] != filepath.Clean(dirA) || dirs[1] != filepath.Clean(dirB) || dirs[2] != filepath.Clean(dirA) || dirs[3] != filepath.Clean(dirB) {
		t.Fatalf("unexpected normalization order: %v", dirs)
	}
	if len(warnings) != 2 {
		t.Fatalf("expected warnings for missing path and file, got %d (%v)", len(warnings), warnings)
	}
}

func TestSeedStartupWorkspaces_REQ_WORKSPACE_START_DIRS(t *testing.T) {
	t.Parallel()

	temp := t.TempDir()
	statePath := filepath.Join(temp, "state.json")
	dirA := filepath.Join(temp, "projectA")
	dirB := filepath.Join(temp, "projectB")
	if err := os.MkdirAll(dirA, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(dirB, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	g := NewGoful(statePath)
	if !SeedStartupWorkspaces(g, []string{dirA, dirB}, false) {
		t.Fatalf("expected seeding to occur")
	}

	ws := g.Workspace()
	if len(ws.Dirs) != 2 {
		t.Fatalf("expected 2 directory windows, got %d", len(ws.Dirs))
	}
	if got := filepath.Clean(ws.Dirs[0].Path); got != filepath.Clean(dirA) {
		t.Fatalf("window 1 mismatch: %s", got)
	}
	if got := filepath.Clean(ws.Dirs[1].Path); got != filepath.Clean(dirB) {
		t.Fatalf("window 2 mismatch: %s", got)
	}

	if SeedStartupWorkspaces(g, nil, false) {
		t.Fatalf("expected empty input to skip seeding")
	}
}
