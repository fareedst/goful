package filer

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func newTestDirectory(t *testing.T, path string) *Directory {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	dir := NewDirectory(0, 0, 80, 20)
	dir.Chdir(path)
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	return dir
}

func namesOf(dir *Directory) []string {
	names := make([]string, 0, dir.Len())
	for _, item := range dir.List() {
		names = append(names, item.Name())
	}
	sort.Strings(names)
	return names
}

func TestFlowOpenDirectory_REQ_INTEGRATION_FLOWS(t *testing.T) {
	// [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_INTEGRATION_FLOWS]
	tmp := t.TempDir()
	files := []string{"alpha.txt", "bravo.txt"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte(name), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	dir := newTestDirectory(t, tmp)
	ws := NewWorkspace(0, 0, 80, 20, "test")
	ws.Dirs = []*Directory{dir}
	ws.SetFocus(0)

	got := namesOf(dir)
	for _, want := range files {
		if !contains(got, want) {
			t.Fatalf("open directory missing %s in %v", want, got)
		}
	}
}

func TestFlowNavigateRename_REQ_INTEGRATION_FLOWS(t *testing.T) {
	// [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_INTEGRATION_FLOWS]
	tmp := t.TempDir()
	subdir := filepath.Join(tmp, "nested")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	orig := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(orig, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	dir := newTestDirectory(t, tmp)
	dir.SetCursorByName("nested")
	dir.EnterDir()
	if dir.Path != filepath.Clean(subdir) {
		t.Fatalf("EnterDir path = %s, want %s", dir.Path, subdir)
	}

	newName := filepath.Join(tmp, "file-renamed.txt")
	if err := os.Rename(orig, newName); err != nil {
		t.Fatalf("rename: %v", err)
	}
	dir.Chdir(tmp)
	dir.reload()
	if idx := dir.IndexByName("file-renamed.txt"); idx == dir.Lower() {
		t.Fatalf("rename result not visible after reload")
	}
}

func TestFlowDelete_REQ_INTEGRATION_FLOWS(t *testing.T) {
	// [REQ:INTEGRATION_FLOWS] [ARCH:TEST_STRATEGY_INTEGRATION] [IMPL:TEST_INTEGRATION_FLOWS]
	tmp := t.TempDir()
	target := filepath.Join(tmp, "remove.me")
	if err := os.WriteFile(target, []byte("bye"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	dir := newTestDirectory(t, tmp)
	if err := os.Remove(target); err != nil {
		t.Fatalf("remove: %v", err)
	}
	dir.reload()
	if idx := dir.IndexByName("remove.me"); idx != dir.Lower() {
		t.Fatalf("deleted entry still present with index %d", idx)
	}
}

func contains(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
